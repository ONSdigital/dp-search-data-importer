package steps

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	componenttest "github.com/ONSdigital/dp-component-test"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-data-importer/config"
	"github.com/ONSdigital/dp-search-data-importer/service"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/maxcnunes/httpfake"
)

const (
	WaitEventTimeout = 5 * time.Second // maximum time that the component test consumer will wait for a kafka event
)

var (
	BuildTime = "1625046891"
	GitCommit = "7434fe334d9f51b7239f978094ea29d10ac33b16"
	Version   = ""
)

type Component struct {
	componenttest.ErrorFeature
	ElasticSearchAPI *httpfake.HTTPFake  // ElasticSearch API mock at HTTP level
	KafkaConsumer    *kafkatest.Consumer // Mock for service kafka consumer
	errorChan        chan error
	svc              *service.Service
	cfg              *config.Config
	wg               *sync.WaitGroup
	signals          chan os.Signal
	waitEventTimeout time.Duration
	ctx              context.Context
}

func NewComponent(t *testing.T) *Component {
	c := &Component{
		ElasticSearchAPI: httpfake.New(
			httpfake.WithTesting(t),
		),
		errorChan:        make(chan error),
		waitEventTimeout: WaitEventTimeout,
		wg:               &sync.WaitGroup{},
		ctx:              context.Background(),
	}
	service.GetKafkaConsumer = c.GetKafkaConsumer
	return c
}

// initService initialises the server, the mocks and waits for the dependencies to be ready
func (c *Component) initService(ctx context.Context) error {
	// register interrupt signals
	c.signals = make(chan os.Signal, 1)
	signal.Notify(c.signals, syscall.SIGINT, syscall.SIGTERM)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	cfg.HealthCheckInterval = time.Second
	cfg.ElasticSearchAPIURL = c.ElasticSearchAPI.ResolveURL("")

	log.Info(ctx, "config used by component tests", log.Data{"cfg": cfg})

	// Create service and initialise it
	c.svc = service.New()
	if err = c.svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("unexpected service Init error in NewComponent: %w", err)
	}

	c.cfg = cfg

	return nil
}

func (c *Component) startService(ctx context.Context) error {
	if err := c.svc.Start(ctx, c.errorChan); err != nil {
		return fmt.Errorf("unexpected error while starting service: %w", err)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		// blocks until an os interrupt or a fatal error occurs
		select {
		case err := <-c.errorChan:
			if errClose := c.svc.Close(ctx); errClose != nil {
				log.Warn(ctx, "error closing server during error handing", log.Data{"close_error": errClose})
			}
			panic(fmt.Errorf("unexpected error received from errorChan: %w", err))
		case sig := <-c.signals:
			log.Info(ctx, "os signal received", log.Data{"signal": sig})
		}

		if err := c.svc.Close(ctx); err != nil {
			panic(fmt.Errorf("unexpected error during service graceful shutdown: %w", err))
		}
	}()

	return nil
}

// Close kills the application under test and waits for it to complete the graceful shutdown, or timeout
func (c *Component) Close() {
	// kill application
	c.signals <- os.Interrupt

	// wait for graceful shutdown to finish (or timeout)
	c.wg.Wait()
}

// Reset re-initialises the service under test and the api mocks.
// Note that the service under test should not be started yet
// to prevent race conditions if it tries to call un-initialised dependencies (steps)
func (c *Component) Reset() error {
	if err := c.initService(c.ctx); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}

	c.ElasticSearchAPI.Reset()

	return nil
}

// GetKafkaConsumer creates a new kafkatest consumer and stores it to the caller Component struct
// It returns the mock, so it can be used by the service under test.
// If there is any error creating the mock, it is also returned to the service.
func (c *Component) GetKafkaConsumer(ctx context.Context, cfg *config.Kafka) (kafka.IConsumerGroup, error) {
	var err error
	c.KafkaConsumer, err = kafkatest.NewConsumer(
		c.ctx,
		&kafka.ConsumerGroupConfig{
			BrokerAddrs:       cfg.Addr,
			Topic:             cfg.PublishedContentTopic,
			GroupName:         cfg.PublishedContentGroup,
			MinBrokersHealthy: &cfg.ConsumerMinBrokersHealthy,
			KafkaVersion:      &cfg.Version,
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafkatest consumer: %w", err)
	}
	return c.KafkaConsumer.Mock, nil
}
