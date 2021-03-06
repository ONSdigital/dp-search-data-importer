package service

import (
	"context"
	"net/http"

	dpES "github.com/ONSdigital/dp-elasticsearch/v3"
	dpESClient "github.com/ONSdigital/dp-elasticsearch/v3/client"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpkafka "github.com/ONSdigital/dp-kafka/v2"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-search-data-importer/config"
	"github.com/ONSdigital/log.go/v2/log"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	HealthCheck   bool
	KafkaConsumer bool
	ElasticSearch bool
	Init          Initialiser
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck:   false,
		KafkaConsumer: false,
		ElasticSearch: false,
		Init:          initialiser,
	}
}

// Init implements the Initialiser interface to initialise dependencies
type Init struct{}

// GetHTTPServer creates an http server and sets the Server flag to true
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router)
	return s
}

// GetKafkaConsumer creates a Kafka consumer and sets the consumer flag to true
func (e *ExternalServiceList) GetKafkaConsumer(ctx context.Context, cfg *config.Config) (dpkafka.IConsumerGroup, error) {
	consumer, err := e.Init.DoGetKafkaConsumer(ctx, cfg)
	if err != nil {
		return nil, err
	}
	e.KafkaConsumer = true
	return consumer, nil
}

// GetElasticSearchClient creates a ElasticSearchClient and sets the ElasticSearchClient flag to true
func (e *ExternalServiceList) GetElasticSearchClient(ctx context.Context, cfg *config.Config) (dpESClient.Client, error) {
	esClient, err := e.Init.DoGetElasticSearchClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	e.ElasticSearch = true
	return esClient, nil
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// DoGetKafkaConsumer returns a Kafka Consumer group
func (e *Init) DoGetKafkaConsumer(ctx context.Context, cfg *config.Config) (dpkafka.IConsumerGroup, error) {
	cgChannels := dpkafka.CreateConsumerGroupChannels(1)

	kafkaOffset := dpkafka.OffsetNewest
	if cfg.KafkaOffsetOldest {
		kafkaOffset = dpkafka.OffsetOldest
	}

	cConfig := &dpkafka.ConsumerGroupConfig{
		KafkaVersion: &cfg.KafkaVersion,
		Offset:       &kafkaOffset,
	}

	// KafkaTLSProtocol informs service to use TLS protocol for kafka
	if cfg.KafkaSecProtocol == config.KafkaTLSProtocol {
		cConfig.SecurityConfig = dpkafka.GetSecurityConfig(
			cfg.KafkaSecCACerts,
			cfg.KafkaSecClientCert,
			cfg.KafkaSecClientKey,
			cfg.KafkaSecSkipVerify,
		)
	}

	kafkaConsumer, err := dpkafka.NewConsumerGroup(
		ctx,
		cfg.KafkaAddr,
		cfg.PublishedContentTopic,
		cfg.PublishedContentGroup,
		cgChannels,
		cConfig,
	)
	if err != nil {
		return nil, err
	}

	return kafkaConsumer, nil
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

// DoGetElasticSearchClient returns a Elastic Search Client
func (e *Init) DoGetElasticSearchClient(ctx context.Context, cfg *config.Config) (dpESClient.Client, error) {
	esConfig := dpESClient.Config{
		ClientLib: dpESClient.GoElasticV710,
		Address:   cfg.ElasticSearchAPIURL,
	}
	esConfig.Transport = awsSigner
	esClient, esClientErr := dpES.NewClient(esConfig)
	if esClientErr != nil {
		log.Error(ctx, "Failed to create dp-elasticsearch client", esClientErr)
		return nil, esClientErr
	}
	return esClient, nil
}
