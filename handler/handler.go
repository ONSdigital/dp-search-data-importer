package handler

import (
	"context"
	"encoding/json"
	"fmt"

	dpelasticsearch "github.com/ONSdigital/dp-elasticsearch/v3/client"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-data-importer/config"
	"github.com/ONSdigital/dp-search-data-importer/models"
	"github.com/ONSdigital/dp-search-data-importer/schema"
	"github.com/ONSdigital/dp-search-data-importer/transform"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	esDestIndex  = "ons"
	esRespCreate = "create"
)

// BatchHandler handles batches of SearchDataImportModel events that contain CSV row data.
type BatchHandler struct {
	esClient dpelasticsearch.Client
	esURL    string
}

// NewBatchHandler returns a BatchHandler.
func NewBatchHandler(esClient dpelasticsearch.Client, cfg *config.Config) *BatchHandler {
	return &BatchHandler{
		esClient: esClient,
		esURL:    cfg.ElasticSearchAPIURL,
	}
}

func (h *BatchHandler) Publish(ctx context.Context, batch []kafka.Message) error {
	// no events received. Nothing more to do (this scenario should not happen)
	if len(batch) == 0 {
		log.Info(ctx, "there are no events to handle")
		return nil
	}

	// unmarshal all events in batch
	events := make([]*models.SearchDataImport, len(batch))
	for i, msg := range batch {
		e := &models.SearchDataImport{}
		s := schema.SearchDataImportEvent

		if err := s.Unmarshal(msg.GetData(), e); err != nil {
			return &Error{
				err: fmt.Errorf("failed to unmarshal event: %w", err),
				logData: map[string]interface{}{
					"msg_data": string(msg.GetData()),
				},
			}
		}

		events[i] = e
	}

	// Summarise what we received (avoid dumping entire structs)
	summaries := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		if e == nil {
			summaries = append(summaries, map[string]interface{}{"error": "nil event"})
			continue
		}
		summaries = append(summaries, map[string]interface{}{
			"uri":          e.URI,
			"search_index": e.SearchIndex,
			"trace_id":     e.TraceID,
		})
	}
	log.Info(ctx, "batch of events received", log.Data{
		"count":  len(events),
		"events": summaries,
	})

	// send batch to elasticsearch
	err := h.sendToES(ctx, events)
	if err != nil {
		log.Error(ctx, "failed to send event to Elastic Search", err)
		return err
	}

	return nil
}

// Preparing the payload and sending bulk events to elastic search.
func (h *BatchHandler) sendToES(ctx context.Context, events []*models.SearchDataImport) error {

	log.Info(ctx, "bulk events into ES starts")
	target := len(events)

	var bulkupsert []byte
	for _, event := range events {
		if event.UID == "" {
			log.Info(ctx, "no uid for inbound kafka event, no transformation possible")
			continue // break here
		}

		upsertBulkRequestBody, err := prepareEventForBulkUpsertRequestBody(ctx, event)
		if err != nil {
			log.Error(ctx, "error in preparing the bulk for upsert", err, log.Data{
				"event": *event,
			})
			continue
		}
		bulkupsert = append(bulkupsert, upsertBulkRequestBody...)
	}

	jsonUpsertResponse, err := h.esClient.BulkUpdate(ctx, esDestIndex, h.esURL, bulkupsert)
	if err != nil {
		if jsonUpsertResponse == nil {
			log.Error(ctx, "server error while upserting the event", err)
			return err
		}
		log.Warn(ctx, "error in response from elasticsearch while upserting the event", log.FormatErrors([]error{err}))
	}

	var bulkRes models.EsBulkResponse
	if err := json.Unmarshal(jsonUpsertResponse, &bulkRes); err != nil {
		log.Error(ctx, "error unmarshaling json", err)
		return err
	}
	if bulkRes.Errors {
		for _, resUpsertItem := range bulkRes.Items {
			if resUpsertItem[esRespCreate].Status == 409 {
				continue
			} else {
				log.Error(ctx, "error upserting doc to ES", err,
					log.Data{
						"response.uid:":   resUpsertItem[esRespCreate].ID,
						"response status": resUpsertItem[esRespCreate].Status,
					})
				target--
				continue
			}
		}
	}

	log.Info(ctx, "documents bulk uploaded to elasticsearch", log.Data{
		"documents_received": len(events),
		"documents_inserted": target,
	})
	return nil
}

// Preparing the payload to be inserted into the elastic search.
func prepareEventForBulkUpsertRequestBody(ctx context.Context, sdModel *models.SearchDataImport) (bulkbody []byte, err error) {

	uid := sdModel.UID
	t := transform.NewTransformer()
	esModel := t.TransformEventModelToEsModel(sdModel)

	if esModel != nil {
		b, err := json.Marshal(esModel)
		if err != nil {
			log.Error(ctx, "error marshal to json while preparing bulk request", err)
			return nil, err
		}

		bulkbody = append(bulkbody, []byte("{ \""+"update"+"\": { \"_id\": \""+uid+"\" } }\n")...)
		bulkbody = append(bulkbody, []byte("{")...)
		bulkbody = append(bulkbody, []byte("\"doc\":")...)
		bulkbody = append(bulkbody, b...)
		bulkbody = append(bulkbody, []byte(",\"doc_as_upsert\": true")...)
		bulkbody = append(bulkbody, []byte("}")...)
		bulkbody = append(bulkbody, []byte("\n")...)
	}
	return bulkbody, nil
}

func (h *BatchHandler) Delete(ctx context.Context, batch []kafka.Message) error {
	if len(batch) == 0 {
		log.Info(ctx, "no delete events to handle")
		return nil
	}

	// Unmarshal all delete events
	events := make([]*models.DeleteEvent, len(batch))
	s := schema.SearchContentDeletedEvent
	for i, msg := range batch {
		e := &models.DeleteEvent{}
		if err := s.Unmarshal(msg.GetData(), e); err != nil {
			return &Error{
				err:     fmt.Errorf("failed to unmarshal event: %w", err),
				logData: map[string]interface{}{"msg_data": string(msg.GetData())},
			}
		}
		events[i] = e
	}

	// Summarise what we received (avoid dumping entire structs)
	summaries := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		if e == nil {
			summaries = append(summaries, map[string]interface{}{"error": "nil event"})
			continue
		}
		summaries = append(summaries, map[string]interface{}{
			"uri":          e.URI,
			"search_index": e.SearchIndex,
			"trace_id":     e.TraceID,
		})
	}
	log.Info(ctx, "batch of delete events received", log.Data{
		"count":  len(events),
		"events": summaries,
	})

	// Build bulk delete body
	bulkBody, err := buildBulkDeleteBody(events)
	if err != nil {
		log.Error(ctx, "failed to build bulk delete body", err)
		return err
	}
	if len(bulkBody) == 0 {
		log.Info(ctx, "no valid delete actions to send (empty bulk body)")
		return nil
	}

	// Send bulk request to Elasticsearch
	jsonResp, err := h.esClient.BulkUpdate(ctx, esDestIndex, h.esURL, bulkBody)
	if err != nil {
		if jsonResp == nil {
			log.Error(ctx, "server error while sending bulk delete", err)
			return err
		}
		log.Warn(ctx, "error in response from elasticsearch while sending bulk delete", log.FormatErrors([]error{err}))
	}

	var bulkRes models.EsBulkResponse
	if err := json.Unmarshal(jsonResp, &bulkRes); err != nil {
		log.Error(ctx, "error unmarshalling bulk delete response", err)
		return err
	}

	total := len(events)
	success, failures := 0, 0
	for _, item := range bulkRes.Items {
		res := item["delete"]
		switch res.Status {
		case 200, 404:
			success++
		default:
			failures++
			log.Error(ctx, "error deleting doc in ES bulk", fmt.Errorf("status %d", res.Status),
				log.Data{"response.id": res.ID, "response.status": res.Status, "response.error": res.Error})
		}
	}

	log.Info(ctx, "bulk deletes sent to elasticsearch", log.Data{
		"documents_received": total,
		"documents_success":  success,
		"documents_failed":   failures,
	})

	if failures > 0 {
		return fmt.Errorf("one or more deletes failed: %d of %d", failures, total)
	}
	return nil
}

// Build the delete payload to be sent to elastic search.
func buildBulkDeleteBody(events []*models.DeleteEvent) ([]byte, error) {
	type del struct {
		ID string `json:"_id"`
	}
	type action struct {
		Delete del `json:"delete"`
	}

	avgLineSize := 60
	body := make([]byte, 0, len(events)*avgLineSize)

	for _, e := range events {
		if e == nil || e.URI == "" {
			continue
		}
		line, err := json.Marshal(action{Delete: del{ID: e.URI}})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal bulk delete line for id %q: %w", e.URI, err)
		}
		body = append(body, line...)
		body = append(body, '\n')
	}
	return body, nil
}
