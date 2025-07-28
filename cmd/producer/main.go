package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-data-importer/config"
	"github.com/ONSdigital/dp-search-data-importer/models"
	"github.com/ONSdigital/dp-search-data-importer/schema"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const serviceName = "dp-search-data-importer"

const (
	ZebedeeDataType = "legacy"
	DatasetDataType = "datasets"
)

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeImport  EventType = iota
	EventTypeDelete
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(ctx, "error getting config", err)
		os.Exit(1)
	}

	importKafkaProducer, deleteKafkaProducer := createKafkaProducers(ctx, cfg)

	time.Sleep(500 * time.Millisecond)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		switch scanEventType(scanner) {
		case EventTypeImport:
			e := scanImportEvent(scanner)
			sendImportEvent(ctx, e, importKafkaProducer)
		case EventTypeDelete:
			e := scanDeleteEvent(scanner)
			sendDeleteEvent(ctx, e, deleteKafkaProducer)
		default:
			log.Fatal(ctx, "unknown event type", errors.New("unknown event type"))
			os.Exit(1)
		}
	}
}

func createKafkaProducers(ctx context.Context, cfg *config.Config) (importKafkaProducer, deleteKafkaProducer *kafka.Producer) {
	importProducerConfig := &kafka.ProducerConfig{
		BrokerAddrs:     cfg.Kafka.Addr,
		Topic:           cfg.Kafka.PublishedContentTopic,
		KafkaVersion:    &cfg.Kafka.Version,
		MaxMessageBytes: &cfg.Kafka.MaxBytes,
	}
	deleteProducerConfig := &kafka.ProducerConfig{
		BrokerAddrs:     cfg.Kafka.Addr,
		Topic:           cfg.Kafka.DeletedContentTopic,
		KafkaVersion:    &cfg.Kafka.Version,
		MaxMessageBytes: &cfg.Kafka.MaxBytes,
	}
	if cfg.Kafka.SecProtocol == config.KafkaTLSProtocol {
		importProducerConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.Kafka.SecCACerts,
			cfg.Kafka.SecClientCert,
			cfg.Kafka.SecClientKey,
			cfg.Kafka.SecSkipVerify,
		)
		deleteProducerConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.Kafka.SecCACerts,
			cfg.Kafka.SecClientCert,
			cfg.Kafka.SecClientKey,
			cfg.Kafka.SecSkipVerify,
		)
	}
	importKafkaProducer, err := kafka.NewProducer(ctx, importProducerConfig)
	if err != nil {
		log.Fatal(ctx, "fatal error trying to create kafka producer", err, log.Data{"topic": cfg.Kafka.PublishedContentTopic})
		os.Exit(1)
	}
	deleteKafkaProducer, err = kafka.NewProducer(ctx, deleteProducerConfig)
	if err != nil {
		log.Fatal(ctx, "fatal error trying to create kafka producer", err, log.Data{"topic": cfg.Kafka.PublishedContentTopic})
		os.Exit(1)
	}

	// kafka error logging go-routines
	importKafkaProducer.LogErrors(ctx)
	deleteKafkaProducer.LogErrors(ctx)

	return importKafkaProducer, deleteKafkaProducer
}

func scanEventType(scanner *bufio.Scanner) EventType {
	fmt.Println("Chose event type (i[mport]/d[elete])")
	fmt.Printf("$ ")
	scanner.Scan()
	uid := strings.ToLower(scanner.Text())
	switch uid {
	case "i", "import":
		return EventTypeImport
	case "d", "delete":
		return EventTypeDelete
	default:
		return EventTypeUnknown
	}
}

// scanImportEvent creates a searchDataImport event according to the user input
func scanImportEvent(scanner *bufio.Scanner) *models.SearchDataImport {
	fmt.Println("--- [Send Kafka PublishedContent] ---")

	fmt.Println("Type the UID")
	fmt.Printf("$ ")
	scanner.Scan()
	uid := scanner.Text()

	fmt.Printf("Type the Data Type (%s or %s)\n", DatasetDataType, ZebedeeDataType)
	fmt.Printf("$ ")
	scanner.Scan()
	dataType := scanner.Text()

	fmt.Println("Type the JobID")
	fmt.Printf("$ ")
	scanner.Scan()
	jobID := scanner.Text()

	fmt.Println("Type the CDID")
	fmt.Printf("$ ")
	scanner.Scan()
	cdid := scanner.Text()

	fmt.Println("Type the DatasetID")
	fmt.Printf("$ ")
	scanner.Scan()
	datasetID := scanner.Text()

	fmt.Println("Type the Summary")
	fmt.Printf("$ ")
	scanner.Scan()
	summary := scanner.Text()

	fmt.Println("Type the title")
	fmt.Printf("$ ")
	scanner.Scan()
	title := scanner.Text()

	searchDataImport := &models.SearchDataImport{
		UID:         uid,
		DataType:    dataType,
		JobID:       jobID,
		SearchIndex: "ONS",
		CDID:        cdid,
		DatasetID:   datasetID,
		Keywords:    []string{"keyword1", "keyword2"},
		Summary:     summary,
		ReleaseDate: "2017-09-07",
		Title:       title,
		TraceID:     "2effer334d",
	}

	if dataType == DatasetDataType {
		popType := models.PopulationType{}
		dimensions := []models.Dimension{}

		fmt.Println("Type the population type Name")
		fmt.Printf("$ ")
		scanner.Scan()
		popType.Name = scanner.Text()

		fmt.Println("Type the population type Label")
		fmt.Printf("$ ")
		scanner.Scan()
		popType.Label = scanner.Text()

		for {
			fmt.Println("Add dimension? [Yy] to confirm")
			fmt.Printf("$ ")
			scanner.Scan()
			if strings.ToLower(scanner.Text()) != "y" {
				break
			}

			dim := models.Dimension{}

			fmt.Println("Type the dimension Name")
			fmt.Printf("$ ")
			scanner.Scan()
			dim.Name = scanner.Text()

			fmt.Println("Type the dimension Raw Label")
			fmt.Printf("$ ")
			scanner.Scan()
			dim.RawLabel = scanner.Text()

			fmt.Println("Type the dimension Label")
			fmt.Printf("$ ")
			scanner.Scan()
			dim.Label = scanner.Text()

			dimensions = append(dimensions, dim)
		}

		searchDataImport.PopulationType = popType
		searchDataImport.Dimensions = dimensions
	}

	return searchDataImport
}

func sendImportEvent(ctx context.Context, e *models.SearchDataImport, importKafkaProducer *kafka.Producer) {
	log.Info(ctx, "sending search data import event", log.Data{"searchDataImportEvent": e})
	eventData, err := schema.SearchDataImportEvent.Marshal(e)
	if err != nil {
		log.Fatal(ctx, "search data import event error", err)
		os.Exit(1)
	}
	// Send bytes to Output channel, after calling Initialise just in case it is not initialised.
	err = importKafkaProducer.Initialise(ctx)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka producer", err)
		os.Exit(1)
	}
	importKafkaProducer.Channels().Output <- eventData
}

func scanDeleteEvent(scanner *bufio.Scanner) *models.DeleteEvent {
	fmt.Println("--- [Send Kafka Delete Event] ---")

	fmt.Println("Type the URI to be deleted")
	fmt.Printf("$ ")
	scanner.Scan()
	uri := scanner.Text()

	fmt.Println("Type the search index (usually `ons`)")
	fmt.Printf("$ ")
	scanner.Scan()
	index := scanner.Text()

	traceID := uuid.NewString()

	deleteEvent := models.DeleteEvent{
		URI:         uri,
		SearchIndex: index,
		TraceID:     traceID,
	}

	return &deleteEvent
}

func sendDeleteEvent(ctx context.Context, e *models.DeleteEvent, deleteKafkaProducer *kafka.Producer) {
	log.Info(ctx, "sending delete event", log.Data{"deleteEvent": e})
	eventData, err := schema.SearchContentDeletedEvent.Marshal(e)
	if err != nil {
		log.Fatal(ctx, "search data import event error", err)
		os.Exit(1)
	}
	// Send bytes to Output channel, after calling Initialise just in case it is not initialised.
	err = deleteKafkaProducer.Initialise(ctx)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka producer", err)
		os.Exit(1)
	}
	deleteKafkaProducer.Channels().Output <- eventData
}
