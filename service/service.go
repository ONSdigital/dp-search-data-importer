package service

import (
	"context"
	"fmt"

	dpelasticsearch "github.com/ONSdigital/dp-elasticsearch/v3/client"
	dpkafka "github.com/ONSdigital/dp-kafka/v5"
	"github.com/ONSdigital/dp-search-data-importer/config"
	"github.com/ONSdigital/dp-search-data-importer/handler"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the event handler service
type Service struct {
	Cfg             *config.Config
	Server          HTTPServer
	HealthCheck     HealthChecker
	PublishConsumer dpkafka.IConsumerGroup
	DeleteConsumer  dpkafka.IConsumerGroup
	EsCli           dpelasticsearch.Client
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Init(ctx context.Context, cfg *config.Config, buildTime, gitCommit, version string) error {
	var err error

	if cfg == nil {
		return errors.New("nil config passed to service init")
	}
	svc.Cfg = cfg

	// Get ElasticSearch client
	svc.EsCli, err = GetElasticSearchClient(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialise elastic search client: %w", err)
	}

	// Create publish consumer
	if svc.PublishConsumer, err = GetKafkaConsumer(ctx, cfg.Kafka, cfg.Kafka.PublishedContentTopic, cfg.Kafka.PublishedContentGroup); err != nil {
		return fmt.Errorf("failed to create kafka publish consumer: %w", err)
	}

	// Create delete consumer
	if svc.DeleteConsumer, err = GetKafkaConsumer(ctx, cfg.Kafka, cfg.Kafka.DeletedContentTopic, cfg.Kafka.PublishedContentGroup); err != nil {
		return fmt.Errorf("failed to create kafka delete consumer: %w", err)
	}

	// Create publish batch handler and register it to the kafka publish consumer
	publishHandler := handler.NewBatchHandler(svc.EsCli, cfg)
	if err := svc.PublishConsumer.RegisterBatchHandler(ctx, publishHandler.Publish); err != nil {
		return fmt.Errorf("failed to register batch publish handler: %w", err)
	}

	// Create delete batch handler and register it to the kafka delete consumer
	deleteHandler := handler.NewBatchHandler(svc.EsCli, cfg)
	if err := svc.DeleteConsumer.RegisterBatchHandler(ctx, deleteHandler.Delete); err != nil {
		return fmt.Errorf("failed to register batch delete handler: %w", err)
	}

	// Get HealthCheck
	svc.HealthCheck, err = GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return fmt.Errorf("could not instantiate healthcheck: %w", err)
	}

	if err := svc.registerCheckers(ctx); err != nil {
		return fmt.Errorf("unable to register checkers: %w", err)
	}

	// Get HTTP Server with collectionID checkHeader middleware
	r := mux.NewRouter()
	r.StrictSlash(true).Path("/health").HandlerFunc(svc.HealthCheck.Handler)
	svc.Server = GetHTTPServer(cfg.BindAddr, r)

	return nil
}

// Start the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) error {
	log.Info(ctx, "starting service")

	// Kafka error logging go-routine
	svc.PublishConsumer.LogErrors(ctx)
	svc.DeleteConsumer.LogErrors(ctx)

	// If start/stop on health updates is disabled, start consuming as soon as possible
	if !svc.Cfg.StopConsumingOnUnhealthy {
		if err := svc.PublishConsumer.Start(); err != nil {
			return fmt.Errorf("consumer failed to start: %w", err)
		}

		if err := svc.DeleteConsumer.Start(); err != nil {
			return fmt.Errorf("delete consumer failed to start: %w", err)
		}
	}

	// Always start HealthCheck.
	// If start/stop on health updates is enabled,
	// the consumer will start consuming on the first healthy update
	// Otherwise it will start consuming straight away.
	svc.HealthCheck.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failure in http listen and serve: %w", err)
		}
	}()

	return nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Cfg.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	hasShutdownError := false

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.HealthCheck != nil {
			log.Info(ctx, "stopping health checker...")
			svc.HealthCheck.Stop()
			log.Info(ctx, "stopped health checker")
		}

		// If kafka consumer exists, stop listening to it.
		// This will automatically stop the event consumer loops and no more messages will be processed.
		// Note that any remaining batch will be consumed and 'StopAndWait' will wait until this happens.
		// The kafka consumer will be closed after the service shuts down.
		//nolint
		if svc.PublishConsumer != nil {
			log.Info(ctx, "stopping kafka consumer listener...")
			if err := svc.PublishConsumer.StopAndWait(); err != nil {
				log.Error(ctx, "error stopping kafka consumer listener", err)
				hasShutdownError = true
			} else {
				log.Info(ctx, "stopped kafka consumer listener")
			}
		}

		if svc.DeleteConsumer != nil {
			log.Info(ctx, "stopping kafka delete consumer listener...")
			if err := svc.DeleteConsumer.StopAndWait(); err != nil {
				log.Error(ctx, "error stopping kafka delete consumer listener", err)
				hasShutdownError = true
			} else {
				log.Info(ctx, "stopped kafka delete consumer listener")
			}
		}

		// Shutdown the HTTP server
		if svc.Server != nil {
			log.Info(ctx, "shutting http server down...")
			if err := svc.Server.Shutdown(ctx); err != nil {
				log.Error(ctx, "failed to shutdown http server", err)
				hasShutdownError = true
			} else {
				log.Info(ctx, "shut down http server")
			}
		}

		// If kafka consumer exists, close it.
		//nolint
		if svc.PublishConsumer != nil {
			log.Info(ctx, "closing kafka consumer")
			if err := svc.PublishConsumer.Close(ctx); err != nil {
				log.Error(ctx, "failed to close kafka consumer", err)
				hasShutdownError = true
			} else {
				log.Info(ctx, "closed kafka consumer")
			}
		}

		if svc.DeleteConsumer != nil {
			log.Info(ctx, "closing kafka delete consumer")
			if err := svc.DeleteConsumer.Close(ctx); err != nil {
				log.Error(ctx, "failed to close kafka delete consumer", err)
				hasShutdownError = true
			} else {
				log.Info(ctx, "closed kafka delete consumer")
			}
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}

	// other error
	if hasShutdownError {
		return errors.New("failed to shutdown gracefully")
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers(ctx context.Context) error {
	hasErrors := false

	chkElasticSearch, err := svc.HealthCheck.AddAndGetCheck("Elasticsearch", svc.EsCli.Checker)
	if err != nil {
		log.Error(ctx, "error creating elasticsearch health check", err)
		hasErrors = true
	}

	if _, err = svc.HealthCheck.AddAndGetCheck("Kafka publish consumer", svc.PublishConsumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for Kafka publish consumer", err)
	}

	if _, err = svc.HealthCheck.AddAndGetCheck("Kafka delete consumer", svc.DeleteConsumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for Kafka delete consumer", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}

	if svc.Cfg.StopConsumingOnUnhealthy {
		svc.HealthCheck.Subscribe(svc.PublishConsumer, chkElasticSearch)
		svc.HealthCheck.Subscribe(svc.DeleteConsumer, chkElasticSearch)
	}

	return nil
}
