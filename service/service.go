package service

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-search-data-importer/config"
	"github.com/ONSdigital/dp-search-data-importer/event"
	"github.com/ONSdigital/dp-search-data-importer/handler"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	dpelasticsearch "github.com/ONSdigital/dp-elasticsearch/v3/client"
	dpkafka "github.com/ONSdigital/dp-kafka/v2"
	dpawsauth "github.com/ONSdigital/dp-net/v2/awsauth"
)

// Service contains all the configs, server and clients to run the event handler service
type Service struct {
	server          HTTPServer
	router          *mux.Router
	serviceList     *ExternalServiceList
	healthCheck     HealthChecker
	consumer        dpkafka.IConsumerGroup
	shutdownTimeout time.Duration
}

var awsSigner *dpawsauth.AwsSignerRoundTripper

// Run the service
func Run(ctx context.Context, serviceList *ExternalServiceList, buildTime, gitCommit, version string,
	svcErrors chan error) (*Service, error) {

	log.Info(ctx, "starting dp-search-data-importer service")

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve service configuration")
	}
	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	// Get Kafka consumer
	kafkaConsumer, err := serviceList.GetKafkaConsumer(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka consumer", err)
		return nil, err
	}

	// Initialse AWS signer
	if cfg.SignElasticsearchRequests {
		awsSigner, err = dpawsauth.NewAWSSignerRoundTripper("", "", cfg.AwsRegion, cfg.AwsService)
		if err != nil {
			log.Error(ctx, "failed to create aws v4 signer", err)
			return nil, err
		}
	}

	// Get Elastic Search Client
	elasticSearchClient, err := serviceList.GetElasticSearchClient(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise Elastic Search Client", err)
		return nil, err
	}

	// handle a batch of events.
	batchHandler := handler.NewBatchHandler(elasticSearchClient)
	eventConsumer := event.NewConsumer()

	// Start listening for event messages.
	eventConsumer.Consume(ctx, kafkaConsumer, batchHandler, cfg)

	// Kafka error logging go-routine
	kafkaConsumer.Channels().LogErrors(ctx, "error received from kafka consumer, topic: "+cfg.PublishedContentTopic)

	// Get HealthCheck
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, hc, kafkaConsumer, elasticSearchClient); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	// Get HTTP Server with collectionID checkHeader middleware
	r := mux.NewRouter()
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	hc.Start(ctx)

	// Run the http server in a new go-routine
	s := serviceList.GetHTTPServer(cfg.BindAddr, r)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		server:          s,
		router:          r,
		serviceList:     serviceList,
		healthCheck:     hc,
		consumer:        kafkaConsumer,
		shutdownTimeout: cfg.GracefulShutdownTimeout,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.shutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutdown gracefully closes up
	var gracefulShutdown bool

	go func() {
		defer cancel()
		var hasShutdownError bool

		// stop healthcheck, as it depends on everything else
		if svc.serviceList.HealthCheck {
			svc.healthCheck.Stop()
		}

		// If kafka consumer exists, stop listening to it.
		// This will automatically stop the event consumer loops and no more messages will be processed.
		// The kafka consumer will be closed after the service shuts down.
		if svc.serviceList.KafkaConsumer {
			log.Info(ctx, "stopping kafka consumer listener")
			if err := svc.consumer.StopListeningToConsumer(ctx); err != nil {
				log.Error(ctx, "error stopping kafka consumer listener", err)
				hasShutdownError = true
			}
			log.Info(ctx, "stopped kafka consumer listener")
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// If kafka consumer exists, close it.
		if svc.serviceList.KafkaConsumer {
			log.Info(ctx, "closing kafka consumer")
			if err := svc.consumer.Close(ctx); err != nil {
				log.Error(ctx, "error closing kafka consumer", err)
				hasShutdownError = true
			}
			log.Info(ctx, "closed kafka consumer")
		}

		if !hasShutdownError {
			gracefulShutdown = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	if !gracefulShutdown {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context,
	hc HealthChecker,
	consumer dpkafka.IConsumerGroup,
	esclient dpelasticsearch.Client) (err error) {

	hasErrors := false

	if err = hc.AddCheck("Elasticsearch", esclient.Checker); err != nil {
		log.Error(ctx, "error creating elasticsearch health check", err)
		hasErrors = true
	}

	if err := hc.AddCheck("Kafka consumer", consumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for Kafka", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}
