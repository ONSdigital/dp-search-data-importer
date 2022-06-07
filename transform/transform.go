package transform

import (
	"github.com/ONSdigital/dp-search-data-importer/models"
)

// Transformer provides an interface by which to transform data from one form to another
type Transformer interface {
	TransformEventModelToEsModel(eventModel *models.SearchDataImportModel) *models.EsModel
}

// Transform provides a concrete implementation of the Transformer interface
type Transform struct{}

// NewTransformer returns a concrete implementation of the Transformer interface
func NewTransformer() Transformer {
	return &Transform{}
}

// TransformModelToEsModel transforms a SearchDataImport into its EsModel counterpart
func (t *Transform) TransformEventModelToEsModel(eventModel *models.SearchDataImportModel) *models.EsModel {
	esModels := models.EsModel{
		DataType:        eventModel.DataType,
		URI:             eventModel.URI,
		JobID:           eventModel.JobID,
		SearchIndex:     eventModel.SearchIndex,
		CDID:            eventModel.CDID,
		DatasetID:       eventModel.DatasetID,
		Keywords:        eventModel.Keywords,
		MetaDescription: eventModel.MetaDescription,
		ReleaseDate:     eventModel.ReleaseDate,
		Summary:         eventModel.Summary,
		Title:           eventModel.Title,
		Topics:          eventModel.Topics,
		Cancelled:       eventModel.Cancelled,
		Finalised:       eventModel.Finalised,
		ProvisionalDate: eventModel.ProvisionalDate,
		Published:       eventModel.Published,
	}
	for _, data := range eventModel.DateChanges {
		esModels.DateChanges = append(esModels.DateChanges, models.ReleaseDateChange(data))
	}
	return &esModels
}
