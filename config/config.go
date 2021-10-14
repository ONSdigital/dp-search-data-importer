package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-search-data-importer
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	KafkaAddr                  []string      `envconfig:"KAFKA_ADDR"`
	KafkaVersion               string        `envconfig:"KAFKA_VERSION"`
	KafkaOffsetOldest          bool          `envconfig:"KAFKA_OFFSET_OLDEST"`
	KafkaNumWorkers            int           `envconfig:"KAFKA_NUM_WORKERS"`
	KafkaSecProtocol           string        `envconfig:"KAFKA_SEC_PROTO"`
	KafkaSecCACerts            string        `envconfig:"KAFKA_SEC_CA_CERTS"`
	KafkaSecClientCert         string        `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	KafkaSecClientKey          string        `envconfig:"KAFKA_SEC_CLIENT_KEY"          json:"-"`
	KafkaSecSkipVerify         bool          `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	KafkaTLSProtocolFlag       string        `envconfig:"KAFKA_TLS_PROTOCOL_FLAG"`	
	PublishedContentGroup      string        `envconfig:"KAFKA_PUBLISHED_CONTENT_GROUP"`
	PublishedContentTopic      string        `envconfig:"KAFKA_PUBLISHED_CONTENT_TOPIC"`
	OutputFilePath             string        `envconfig:"OUTPUT_FILE_PATH"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:25900",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		KafkaAddr:                  []string{"localhost:9092"},
		KafkaVersion:               "1.0.2",
		KafkaOffsetOldest:          true,
		KafkaNumWorkers:            1,
		KafkaSecProtocol:           "",
		KafkaSecCACerts:            "",
		KafkaSecClientCert:         "",
		KafkaSecClientKey:          "",
		KafkaSecSkipVerify:         false,
		KafkaTLSProtocolFlag:        "TLS",
		PublishedContentGroup:      "dp-search-data-importer",
		PublishedContentTopic:      "search-data-import",
		OutputFilePath:             "/tmp/search-data-importer.txt",
	}

	return cfg, envconfig.Process("", cfg)
}
