package handler_test

import (
	"context"
	"errors"
	"testing"

	dpelasticsearch "github.com/ONSdigital/dp-elasticsearch/v3/client"
	dpMock "github.com/ONSdigital/dp-elasticsearch/v3/client/mocks"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-data-importer/config"
	"github.com/ONSdigital/dp-search-data-importer/handler"
	"github.com/ONSdigital/dp-search-data-importer/models"
	"github.com/ONSdigital/dp-search-data-importer/schema"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testContext = context.Background()
	indexName   = "ons"
	esDestURL   = "locahost:9999"
	testCfg     = &config.Config{
		ElasticSearchAPIURL: esDestURL,
	}

	expectedEvent1 = &models.SearchDataImport{
		UID:             "uid1",
		DataType:        "anyDataType1",
		JobID:           "",
		SearchIndex:     "ONS",
		CDID:            "",
		DatasetID:       "",
		Keywords:        []string{"anykeyword1"},
		MetaDescription: "",
		Summary:         "",
		ReleaseDate:     "",
		Title:           "anyTitle1",
		TraceID:         "anyTraceID1",
		DateChanges:     []models.ReleaseDateDetails{},
		Cancelled:       false,
		Finalised:       false,
		ProvisionalDate: "",
		Published:       false,
		Survey:          "",
		Language:        "",
		CanonicalTopic:  "",
		PopulationType: models.PopulationType{
			Name:  "pop1",
			Label: "popLbl1",
		},
		Dimensions: []models.Dimension{
			{Name: "dim1", Label: "dimLbl1", RawLabel: "dimRawLbl1"},
		},
	}

	expectedEvent2 = &models.SearchDataImport{
		UID:             "uid2",
		DataType:        "anyDataType2",
		JobID:           "",
		SearchIndex:     "ONS",
		CDID:            "",
		DatasetID:       "",
		Keywords:        []string{"anykeyword2"},
		MetaDescription: "",
		Summary:         "",
		ReleaseDate:     "",
		Title:           "anyTitle2",
		TraceID:         "anyTraceID2",
	}

	testEvents = []*models.SearchDataImport{
		expectedEvent1,
		expectedEvent2,
	}
	mockSuccessESResponseWithNoError = []byte("{\"took\":6,\"errors\":false,\"items\":[{\"create\":{\"_index\":\"ons1637667136829001\",\"_type\":\"_doc\",\"_id\":\"testTitle3\",\"_version\":1,\"result\":\"created\",\"_shards\":{\"total\":2,\"successful\":2,\"failed\":0},\"_seq_no\":0,\"_primary_term\":1,\"status\":201}}]}")
)

func TestPublishHandleWithEventsCreated(t *testing.T) {
	Convey("Given a handler configured with successful es updates for all two events is success", t, func() {
		elasticSearchMock := &dpMock.ClientMock{
			BulkUpdateFunc: func(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error) {
				return mockSuccessESResponseWithNoError, nil
			},
		}

		batchHandler := handler.NewBatchHandler(elasticSearchMock, testCfg)

		Convey("When handle is called with no error", func() {
			err := batchHandler.Publish(testContext, createTestBatch(testEvents))

			Convey("Then the error is nil and it performed upsert action to the expected index", func() {
				So(err, ShouldBeNil)
				So(elasticSearchMock.BulkUpdateCalls(), ShouldHaveLength, 1)
				So(elasticSearchMock.BulkUpdateCalls()[0].IndexName, ShouldEqual, indexName)
				So(elasticSearchMock.BulkUpdateCalls()[0].URL, ShouldEqual, esDestURL)
				So(elasticSearchMock.BulkUpdateCalls()[0].Settings, ShouldNotBeEmpty)
			})
		})
	})
}

func TestPublishHandleWithEventsUpdated(t *testing.T) {
	Convey("Given a handler configured with sucessful es updates for two events with one create error", t, func() {
		elasticSearchMock := &dpMock.ClientMock{
			BulkUpdateFunc: func(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error) {
				return mockSuccessESResponseWithNoError, nil
			},
		}

		batchHandler := handler.NewBatchHandler(elasticSearchMock, testCfg)

		Convey("When handle is called", func() {
			err := batchHandler.Publish(testContext, createTestBatch(testEvents))

			Convey("Then the error is nil and it performed upsert action to the expected index", func() {
				So(err, ShouldBeNil)
				So(elasticSearchMock.BulkUpdateCalls(), ShouldHaveLength, 1)
				So(elasticSearchMock.BulkUpdateCalls()[0].IndexName, ShouldEqual, indexName)
				So(elasticSearchMock.BulkUpdateCalls()[0].URL, ShouldEqual, esDestURL)
				So(elasticSearchMock.BulkUpdateCalls()[0].Settings, ShouldNotBeEmpty)
			})
		})
	})
}

func TestPublishHandleWithInternalServerESResponse(t *testing.T) {
	Convey("Given a handler configured with other failed es create request", t, func() {
		elasticSearchMock := &dpMock.ClientMock{
			BulkUpdateFunc: func(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error) {
				return nil, errors.New("unexpected status code from api")
			},
		}
		batchHandler := handler.NewBatchHandler(elasticSearchMock, testCfg)
		Convey("When handle is called", func() {
			err := batchHandler.Publish(testContext, createTestBatch(testEvents))

			Convey("And the error is not nil while performing upsert action", func() {
				So(err, ShouldResemble, errors.New("unexpected status code from api"))
			})
		})
	})
}

func TestDeleteHandleWithSuccess(t *testing.T) {
	Convey("Given a batch handler with a successful ES delete", t, func() {
		var capturedSearch dpelasticsearch.Search

		mockES := &dpMock.ClientMock{
			DeleteDocumentByQueryFunc: func(ctx context.Context, search dpelasticsearch.Search) error {
				capturedSearch = search
				return nil
			},
		}

		batchHandler := handler.NewBatchHandler(mockES, testCfg)

		deleteEvent := models.DeleteEvent{URI: "/to/delete"}
		msgBytes, err := schema.SearchContentDeletedEvent.Marshal(&deleteEvent)
		So(err, ShouldBeNil)
		msg, err := kafkatest.NewMessage(msgBytes, 1)
		So(err, ShouldBeNil)

		Convey("When Delete is called with a valid delete event", func() {
			err := batchHandler.Delete(testContext, []kafka.Message{msg})

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(mockES.DeleteDocumentByQueryCalls(), ShouldHaveLength, 1)
				So(capturedSearch.Header.Index, ShouldEqual, "ons")
				So(string(capturedSearch.Query), ShouldContainSubstring, "/to/delete")
			})
		})
	})
}

func TestDeleteHandleWithUnmarshalError(t *testing.T) {
	Convey("Given a batch handler with malformed JSON", t, func() {
		mockES := &dpMock.ClientMock{}
		batchHandler := handler.NewBatchHandler(mockES, testCfg)

		badMsg, _ := kafkatest.NewMessage([]byte("bad json"), 1)

		Convey("When Delete is called", func() {
			err := batchHandler.Delete(testContext, []kafka.Message{badMsg})

			Convey("Then an *Error is returned with unmarshal context", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to unmarshal event")

				typedErr, ok := err.(*handler.Error)
				So(ok, ShouldBeTrue)
				So(typedErr.LogData(), ShouldContainKey, "msg_data")
			})
		})
	})
}

func TestDeleteHandleWithDeleteFails(t *testing.T) {
	Convey("Given a batch handler where ES delete fails", t, func() {
		mockES := &dpMock.ClientMock{
			DeleteDocumentByQueryFunc: func(ctx context.Context, search dpelasticsearch.Search) error {
				return errors.New("ES delete failed")
			},
		}

		batchHandler := handler.NewBatchHandler(mockES, testCfg)

		deleteEvent := models.DeleteEvent{URI: "/fail/delete"}
		msgBytes, err := schema.SearchContentDeletedEvent.Marshal(&deleteEvent)
		So(err, ShouldBeNil)
		msg, err := kafkatest.NewMessage(msgBytes, 1)
		So(err, ShouldBeNil)

		Convey("When Delete is called", func() {
			err := batchHandler.Delete(testContext, []kafka.Message{msg})

			Convey("Then an *Error is returned and delete was attempted", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to delete document for uri")
				So(err.Error(), ShouldContainSubstring, "/fail/delete")

				typedErr, ok := err.(*handler.Error)
				So(ok, ShouldBeTrue)
				So(typedErr.LogData(), ShouldContainKey, "uri")

				So(mockES.DeleteDocumentByQueryCalls(), ShouldHaveLength, 1)
			})
		})
	})
}

func createTestBatch(events []*models.SearchDataImport) []kafka.Message {
	batch := make([]kafka.Message, len(events))
	for i, s := range events {
		e, err := schema.SearchDataImportEvent.Marshal(s)
		So(err, ShouldBeNil)
		msg, err := kafkatest.NewMessage(e, int64(i))
		So(err, ShouldBeNil)
		batch[i] = msg
	}
	return batch
}
