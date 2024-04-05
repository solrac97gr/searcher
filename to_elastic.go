package searcher

import (
	"encoding/json"
	"fmt"

	"github.com/solrac97gr/searcher/domain/models"
	"github.com/solrac97gr/searcher/internal/sentinels"
)

// Define the operators as constants for avoid typos and easy editing them later
const (
	queryKey   string = "query"
	boolQuery  string = "bool"
	must       string = "must"
	should     string = "should"
	term       string = "term"
	mustNot    string = "must_not"
	rangeQuery string = "range"
	gt         string = "gt"
	gte        string = "gte"
	lt         string = "lt"
	lte        string = "lte"
	sort       string = "sort"
	order      string = "order"
)

// ToElastic converts criteria to an Elasticsearch query string.
// It takes a validMapEntityName string, a criteria models.Criteria, superFilters array as input.
// It returns a string representing the Elasticsearch query and an error if any.
//
// The function converts the criteria into an Elasticsearch query by applying the specified filters and sorts with a set of SuperFilters in the top of the query that logically ends like (CLIENT_ID="example" AND (THE_QUERY)).
// It checks if the entity has the permitted fields registered in the valid field map.
// If the limit is not specified in the pagination, it defaults to 50.
// The function builds the query using the specified filters and logical operators.
// It handles various operators such as Equals, NotEquals, GreaterThan, LessThan, GreaterAndEqualsThan, and LessAndEqualsThan.
// The resulting query is returned as a JSON string.
func (ca *QueryTranslator) ToElastic(validMapEntityName string, rawCriteria models.Criteria, superFilters []models.SuperFilter) (string, error) {
	// We need to pre-process the criteria adding default values in case of some conditions are matched (check PrepareCriteria())
	criteria := ca.PrepareCriteria(&rawCriteria)

	// Initialize the query map for avoid nil queries
	query := make(map[string]interface{})
	// Initialize the combined query map for avoid nil combined queries
	combinedQuery := make([]map[string]interface{}, 0)

	// Check if the entity has the permitted fields register in the valid field map
	vf, ok := ca.ValidFieldMaps[validMapEntityName]
	if !ok {
		return "", fmt.Errorf("invalid field: %s not valid field registers", validMapEntityName)
	}

	// build the sorts using the helper function BuildSorts
	sorts, err := BuildSorts(criteria.Query.Sorts, vf)
	if err != nil {
		return "", err
	}
	// Assign the sort
	query[sort] = sorts

	// Iterate through the Filters inside of the Query
	for _, filter := range criteria.Query.Filters {
		// Initialize the conditions map for avoid nil condition map
		conditions := make([]map[string]interface{}, 0)
		// Initialize the GreaterThan map for avoid nil GreaterThan map
		gtMaps := make(map[string]models.Condition)
		// Initialize the LessThan map for avoid nil LessThan map
		ltMaps := make(map[string]models.Condition)

		for _, condition := range filter.Conditions {
			// Convert the condition Field to string for easier comparison
			field := condition.Field.String()

			// Check if the field is a valid field for search query
			fieldMetaData, ok := vf.Fields[field]
			if !ok {
				return "", fmt.Errorf("%w:invalid field: %s", sentinels.ErrValidation, field)
			}
			// If the field is marked as analyzed field we add ".raw" to the field name this is because the raw value of the field is storage in the key with name raw
			if fieldMetaData.IsAnalyzed {
				field = field + ".raw"
			}
			// Check if the Field is type Date for convert it from ISO Date string to Unix timestamp
			if fieldMetaData.Type.Equals(models.Date) {
				// Convert the Value (interface{}) to a string
				date, ok := condition.Value.(string)
				if !ok {
					return "", fmt.Errorf("%w:invalid date field: %s", sentinels.ErrValidation, field)
				}
				// Perform the conversion from ISO to Unix timestamp
				newDate, err := ca.dateFormatter.FromISO8601String(date)
				if err != nil {
					return "", fmt.Errorf("%w:invalid date field: %s", sentinels.ErrValidation, field)
				}
				// Assign the new processed value
				condition.Value = newDate
			}
			// Assign the operator from the condition for the switch
			operator := condition.Operator

			// Using a switch case we append the respective conditions to their respective ElasticSearch formatted conditions except for the >,>= and <,<= conditions that are compiled in a map for another processing (determinate if ranges exists) step before to be formatted as ElasticSearch format
			switch operator {
			case models.EqualsOperator:
				conditions = append(conditions, createEqualsCondition(field, condition.Value))
			case models.NotEqualsOperator:
				conditions = append(conditions, createNotEqualsCondition(field, condition.Value))
			case models.GreaterThan:
				gtMaps[field] = condition
			case models.LessThan:
				ltMaps[field] = condition
			case models.GreaterAndEqualsThan:
				gtMaps[field] = condition
			case models.LessAndEqualsThan:
				ltMaps[field] = condition
			}
		}

		// Process the conditions for group the conditions that are has a common field in range conditions
		processRangeConditions(&conditions, &gtMaps, &ltMaps)
		// Add the conditions that are already processed to the combined query
		addFilterConditions(&combinedQuery, filter.Logical, conditions)
	}

	// Build the query after all conditions have been processed into the Elasticsearch format
	buildQuery(superFilters, &query, &combinedQuery, criteria.Query.Logical)

	// Set the pagination parameters for the query
	query["size"] = criteria.Pagination.Limit
	query["from"] = criteria.Pagination.Offset

	// marshal the query for after converting to a string
	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	return string(jsonQuery), nil
}

// BuildSorts builds the sorts for the given query and validate if the sorting fields are valid
func BuildSorts(sorts []models.Sort, vf models.ValidFields) ([]map[string]interface{}, error) {
	buildedSorts := make([]map[string]interface{}, 0)

	for _, srt := range sorts {
		fieldName := srt.Field
		fmd, ok := vf.Fields[fieldName]
		if !ok {
			return nil, fmt.Errorf("%w: %s not valid field registers", sentinels.ErrValidation, fieldName)
		}
		// If the field is analyzed we add the .raw in the name of the field
		// these is the right way for perform and order in analyzed fields in
		// elasticsearch
		if fmd.IsAnalyzed {
			fieldName = fmt.Sprintf("%s.raw", fieldName)
		}

		nSort := map[string]interface{}{
			string(fieldName): map[string]string{
				order: srt.Order.String(),
			},
		}
		buildedSorts = append(buildedSorts, nSort)
	}

	// Add the default sort field (bayonet_tracking_id)
	// we add direct in the buildedSorts for skip the validation of the fields
	// this will permit the field be ordered independently if can be queried
	defaultSort := map[string]interface{}{
		"bayonet_tracking_id": map[string]string{
			order: "asc",
		},
	}
	buildedSorts = append(buildedSorts, defaultSort)

	return buildedSorts, nil
}

// createEqualsCondition helper function for create a equal condition for Elasticsearch
func createEqualsCondition(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		term: map[string]interface{}{
			field: value,
		},
	}
}

// createNotEqualsCondition helper function for creating a not equal condition for Elasticsearch
func createNotEqualsCondition(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		boolQuery: map[string]interface{}{
			mustNot: []map[string]interface{}{
				{
					term: map[string]interface{}{
						field: value,
					},
				},
			},
		},
	}
}

// processRangeConditions take the conditions and the maps of GreaterThan and LessThan conditions for create the ElasticSearch formatted range conditions depending of their type of logical operator and if exists common fields
func processRangeConditions(conditions *[]map[string]interface{}, gtMaps *map[string]models.Condition, ltMaps *map[string]models.Condition) {
	commonKeys := findCommonKeys(*gtMaps, *ltMaps)

	for _, key := range commonKeys {
		*conditions = append(*conditions, createRangeCondition(
			key,
			(*gtMaps)[key].Operator,
			(*ltMaps)[key].Operator,

			(*gtMaps)[key].Value,
			(*ltMaps)[key].Value,
		),
		)
		delete(*gtMaps, key)
		delete(*ltMaps, key)
	}

	for key := range *gtMaps {
		if (*gtMaps)[key].Operator.Equals(models.GreaterAndEqualsThan) {
			*conditions = append(*conditions, createGreaterAndEqualsThanCondition(key, (*gtMaps)[key].Value))
		}
		*conditions = append(*conditions, createGreaterThanCondition(key, (*gtMaps)[key].Value))
	}

	for key := range *ltMaps {
		if (*ltMaps)[key].Operator.Equals(models.LessAndEqualsThan) {
			*conditions = append(*conditions, createLessAndEqualsThanCondition(key, (*ltMaps)[key].Value))
		}
		*conditions = append(*conditions, createLessThanCondition(key, (*ltMaps)[key].Value))
	}
}

// createRangeCondition helper function for creating a range condition for ElasticSearch
func createRangeCondition(key string, gtOperator models.Operator, ltOperator models.Operator, gtValue interface{}, ltValue interface{}) map[string]interface{} {
	gtOp := gt
	ltOp := lt

	if gtOperator.Equals(models.GreaterAndEqualsThan) {
		gtOp = gte
	}

	if ltOperator.Equals(models.LessAndEqualsThan) {
		ltOp = lte
	}

	return map[string]interface{}{
		rangeQuery: map[string]interface{}{
			key: map[string]interface{}{
				gtOp: gtValue,
				ltOp: ltValue,
			},
		},
	}
}

// createGreaterThanCondition helper function for create a formatted Greater Than condition for ElasticSearch
func createGreaterThanCondition(key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		rangeQuery: map[string]interface{}{
			key: map[string]interface{}{
				gt: value,
			},
		},
	}
}

// createLessThanCondition helper function for create a formatted Less Than condition for ElasticSearch
func createLessThanCondition(key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		rangeQuery: map[string]interface{}{
			key: map[string]interface{}{
				lt: value,
			},
		},
	}
}

// createGreaterAndEqualThanCondition helper function for create a formatted Greater and Equal condition for ElasticSearch
func createGreaterAndEqualsThanCondition(key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		rangeQuery: map[string]interface{}{
			key: map[string]interface{}{
				gte: value,
			},
		},
	}
}

// createLessAndEqualsThanCondition helper function for create a formatted Less and equal condition for Elasticsearch
func createLessAndEqualsThanCondition(key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		rangeQuery: map[string]interface{}{
			key: map[string]interface{}{
				lte: value,
			},
		},
	}
}

// addFilterConditions add the conditions to the filter depending on the logical operator set for that filter
//
// - AND = MUST
//
// - OR = SHOULD
func addFilterConditions(combinedQuery *[]map[string]interface{}, logical models.Logical, conditions []map[string]interface{}) {
	if logical.Equals(models.ORLogical) && len(conditions) > 0 {
		*combinedQuery = append(*combinedQuery, map[string]interface{}{
			boolQuery: map[string]interface{}{
				should: conditions,
			},
		})
	} else if logical.Equals(models.ANDLogical) && len(conditions) > 0 {
		*combinedQuery = append(*combinedQuery, map[string]interface{}{
			boolQuery: map[string]interface{}{
				must: conditions,
			},
		})
	}
}

// buildQuery a helper function to build a query adding the top level operator and the super filters to the conditions
//
// - AND = MUST
//
// - OR = SHOULD
//
// The super filters it adds a top level operator extra with always a MUST condition where the super filters are set
func buildQuery(superFilters []models.SuperFilter, query *map[string]interface{}, combinedQuery *[]map[string]interface{}, logical models.Logical) {

	// We build the super filters to the Query
	queryWithSuperFilters := []map[string]interface{}{}
	for _, superFilter := range superFilters {
		queryWithSuperFilters = append(queryWithSuperFilters, map[string]interface{}{
			"term": map[string]interface{}{
				superFilter.Field: superFilter.Value,
			},
		})
	}

	if logical.Equals(models.ANDLogical) {
		if len(*combinedQuery) > 0 {
			queryWithSuperFilters = append(queryWithSuperFilters, map[string]interface{}{
				boolQuery: map[string]interface{}{
					must: *combinedQuery,
				},
			})
		}
		(*query)[queryKey] = map[string]interface{}{
			boolQuery: map[string]interface{}{
				must: queryWithSuperFilters,
			},
		}
	} else if logical.Equals(models.ORLogical) {
		if len(*combinedQuery) > 0 {
			queryWithSuperFilters = append(queryWithSuperFilters, map[string]interface{}{
				boolQuery: map[string]interface{}{
					should: *combinedQuery,
				},
			})
		}
		(*query)[queryKey] = map[string]interface{}{
			boolQuery: map[string]interface{}{
				must: queryWithSuperFilters,
			},
		}
	}

}
