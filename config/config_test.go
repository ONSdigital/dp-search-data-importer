package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		os.Clearenv()
		cfg, err := Get()

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the values should be set to the expected defaults", func() {
				So(cfg.AwsRegion, ShouldEqual, "eu-west-2")
				So(cfg.AwsService, ShouldEqual, "es")
				So(cfg.BatchSize, ShouldEqual, 500)
				So(cfg.BatchWaitTime, ShouldEqual, time.Second*5)
				So(cfg.BindAddr, ShouldEqual, "localhost:25900")
				So(cfg.ElasticSearchAPIURL, ShouldEqual, "http://localhost:11200")
				So(cfg.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(cfg.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.Kafka.Addr, ShouldResemble, []string{"localhost:9092", "localhost:9093", "localhost:9094"})
				So(cfg.Kafka.ConsumerMinBrokersHealthy, ShouldEqual, 1)
				So(cfg.Kafka.MaxBytes, ShouldEqual, 2000000)
				So(cfg.Kafka.NumWorkers, ShouldEqual, 1)
				So(cfg.Kafka.OffsetOldest, ShouldBeTrue)
				So(cfg.Kafka.ProducerMinBrokersHealthy, ShouldEqual, 1)
				So(cfg.Kafka.PublishedContentGroup, ShouldEqual, "dp-search-data-importer")
				So(cfg.Kafka.PublishedContentTopic, ShouldEqual, "search-data-import")
				So(cfg.Kafka.SecCACerts, ShouldEqual, "")
				So(cfg.Kafka.SecClientCert, ShouldEqual, "")
				So(cfg.Kafka.SecClientKey, ShouldEqual, "")
				So(cfg.Kafka.SecProtocol, ShouldEqual, "")
				So(cfg.Kafka.SecSkipVerify, ShouldBeFalse)
				So(cfg.Kafka.Version, ShouldEqual, "1.0.2")
				So(cfg.SignElasticsearchRequests, ShouldEqual, false)
				So(cfg.StopConsumingOnUnhealthy, ShouldBeTrue)
			})

			Convey("Then a second call to config should return the same config", func() {
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}
