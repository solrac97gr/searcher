package searcher

import (
	"fmt"

	"github.com/solrac97gr/searcher/domain/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/solrac97gr/searcher/internal/sentinels"
)

func (ca *QueryTranslator) ToMongo(validMapEntityName string, rawCriteria models.Criteria, superFilters []models.SuperFilter) (models.MongoQuery, error) {
	// We need to pre-process the criteria adding default values in case of some conditions are matched (check PrepareCriteria())
	c := ca.PrepareCriteria(&rawCriteria)

	query := make(map[string]interface{})

	// Add pagination to the query
	query["limit"] = c.Pagination.Limit
	query["offset"] = c.Pagination.Offset

	// Check if the valid fields are correctly registered
	vf, ok := ca.ValidFieldMaps[validMapEntityName]
	if !ok {
		return nil, fmt.Errorf("%w: %s not valid field registers", sentinels.ErrValidation, validMapEntityName)
	}

	// Add filters to the query
	filters := bson.A{}
	for _, filter := range c.Query.Filters {
		conditions := bson.A{}
		for _, condition := range filter.Conditions {
			field := condition.Field.String()

			fieldMetaData, ok := vf.Fields[field]
			if !ok {
				return nil, fmt.Errorf("%w:invalid field: %s", sentinels.ErrValidation, field)
			}
			if fieldMetaData.Type.Equals(models.Date) {
				date, ok := condition.Value.(string)
				if !ok {
					return nil, fmt.Errorf("%w:invalid field: %s", sentinels.ErrValidation, field)
				}
				nDate, err := ca.dateFormatter.FromISO8601String(date)
				if err != nil {
					return nil, fmt.Errorf("%w:invalid field: %s", sentinels.ErrValidation, field)
				}
				condition.Value = nDate
			}

			operator := condition.Operator
			if operator.Equals(models.NotEqualsOperator) {
				operator = "$ne"
			} else if operator.Equals(models.EqualsOperator) {
				operator = "$eq"
			} else if operator.Equals(models.GreaterThan) {
				operator = "$gt"
			} else if operator.Equals(models.LessThan) {
				operator = "$lt"
			} else if operator.Equals(models.LessAndEqualsThan) {
				operator = "$lte"
			} else if operator.Equals(models.GreaterAndEqualsThan) {
				operator = "$gte"
			}
			conditions = append(conditions, bson.M{string(condition.Field): bson.M{operator.String(): condition.Value}})
		}
		filters = append(filters, bson.M{"$" + string(filter.Logical): conditions})
	}
	query["filters"] = bson.M{
		"$and": func() bson.A {
			// We add the super filters to the Top Level Query.
			sf := bson.A{}
			for _, superFilter := range superFilters {
				sf = append(sf, bson.M{
					superFilter.Field: superFilter.Value,
				})
			}
			return sf
		}(),
	}
	if len(filters) > 0 {
		query["filters"].(bson.M)["$and"] = append(query["filters"].(bson.M)["$and"].(bson.A), bson.M{"$" + c.Query.Logical.String(): filters})
	}

	// Add sort to the query
	sort := bson.M{}
	for _, s := range c.Query.Sorts {
		_, ok := vf.Fields[s.Field]
		if !ok {
			return nil, fmt.Errorf("%w:invalid field: %s", sentinels.ErrValidation, s.Field)
		}

		order := 1
		if s.Order.Equals(models.DESCOrder) {
			order = -1
		}
		sort[s.Field] = order
	}

	query["sorts"] = sort

	return query, nil
}
