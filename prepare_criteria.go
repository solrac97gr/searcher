package searcher

import "github.com/solrac97gr/searcher/domain/models"

const (
	// DefaultLogicOperator is the default logic operator
	DefaultLogicOperator = models.ANDLogical
	// DefaultPaginationLimit is the default pagination limit
	DefaultPaginationLimit uint = 50
	// MaximumLimitOffsetSize
	MaximumLimitOffsetSize uint = 10000
)

// PrepareCriteria is a helper function that set all the default values for a criteria independently of to what database will be executed. This will help for ensure consistency between different databases.
func (ca *QueryTranslator) PrepareCriteria(criteria *models.Criteria) *models.Criteria {
	// If the pagination Limit is set to 0 then we set the DefaultValue for the Pagination Limit
	if criteria.Pagination.Limit == 0 {
		criteria.Pagination.Limit = DefaultPaginationLimit
	}
	// If the number of filter is 1 or less we use the default logic Operator
	if criteria.Query.Filters.Len() <= 1 {
		criteria.Query.Logical = DefaultLogicOperator
	}

	preparedFilters := make([]models.Filter, criteria.Query.Filters.Len())
	// If the number of conditions inside of a filter is 1 we use the default logic Operator
	for index, filter := range criteria.Query.Filters {
		if filter.Conditions.Len() == 1 {
			filter.Logical = DefaultLogicOperator
		}
		preparedFilters[index] = filter
	}

	criteria.Query.Filters = preparedFilters

	return criteria
}
