package searcher

import (
	"errors"

	"github.com/solrac97gr/searcher/domain/models"
	"github.com/solrac97gr/searcher/domain/ports"
	"github.com/solrac97gr/searcher/internal/date"
	"github.com/solrac97gr/searcher/internal/date/domain"
)

type QueryTranslator struct {
	dateFormatter  domain.DateFormatter
	ValidFieldMaps map[string]models.ValidFields
}

var _ ports.QueryTranslator = &QueryTranslator{}

func NewQueryTranslator() (*QueryTranslator, error) {
	df, err := date.NewFormatter()
	if err != nil {
		return nil, err
	}

	return &QueryTranslator{
		ValidFieldMaps: make(map[string]models.ValidFields),
		dateFormatter:  df,
	}, nil
}

// AddValidFieldSet Add a new valid field set for a determined entity this only can be set one time every runtime
// if there is another ValidFields set for an entity it will return a error
func (ca *QueryTranslator) AddValidFieldsSet(validFields models.ValidFields) error {
	_, ok := ca.ValidFieldMaps[validFields.EntityName]
	if ok {
		return errors.New("the valid fields already set for this entity")
	}
	ca.ValidFieldMaps[validFields.EntityName] = validFields
	return nil
}

func findCommonKeys(map1, map2 map[string]models.Condition) []string {
	commonKeys := make([]string, 0)

	for key := range map1 {
		if _, ok := map2[key]; ok {
			commonKeys = append(commonKeys, key)
		}
	}

	return commonKeys
}
