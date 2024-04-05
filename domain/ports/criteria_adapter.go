package ports

import (
	"github.com/solrac97gr/searcher/domain/models"
)

// QueryTranslator is an interface that provides methods for converting criteria objects to different query formats.
type QueryTranslator interface {
	// ToMongo converts the given criteria to a MongoDB query.
	// It takes as parameters.
	//
	// - validMapEntityName: The name of the set of valid fields that will be validate.
	//
	// - rawCriteria: A Criteria object to convert to a MongoDB query
	//
	// - superFilters: A list of filters to apply in top-level of the query that only will allow equals operator and will skip validation of valid filters for client.
	//
	// If there is an error during conversion, it returns an error.
	ToMongo(validMapEntityName string, rawCriteria models.Criteria, superFilters []models.SuperFilter) (map[string]interface{}, error)

	// ToElastic converts the given criteria to an Elasticsearch query string.
	// It takes as parameters.
	//
	// - validMapEntityName: The name of the set of valid fields that will be validate.
	//
	// - rawCriteria: A Criteria object to convert to a Elastic query
	//
	// - superFilters: A list of filters to apply in top-level of the query that only will allow equals operator and will skip validation of valid filters for client.
	//
	// If there is an error during conversion, it returns an error.
	ToElastic(validMapEntityName string, criteria models.Criteria, superFilters []models.SuperFilter) (string, error)

	// SetValidFields
	AddValidFieldsSet(validFields models.ValidFields) error
}
