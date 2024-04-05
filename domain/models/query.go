package models

import "fmt"

// Query is the structure that contains the filters and sorts to apply to the search.
type Query struct {
	// Filters is an array of filter to apply to the query
	// - If the filter is empty the query will execute successfully without any further filtering
	Filters Filters `json:"filters"`
	// Sort is an array of sort to apply to the query
	// - If the sort is empty the query will execute successfully without any further sorting (only default database sorting)
	Sorts Sorts `json:"sorts"`
	// Logical is the logic operation to apply to the group of filters
	// - If the len of filters is 1 the logical operator will be "and"
	Logical Logical `json:"logical"`
}

func (q *Query) Validate() error {
	if err := q.Filters.Validate(); err != nil {
		return err
	}
	if err := q.Sorts.Validate(); err != nil {
		return err
	}

	// Validate the logical operator when exist always. Since for 1 filter is optional we must validate if the logical operator is present even if it will be replace in future stages as "AND" like a default operator.
	if q.Logical.String() != "" {
		if err := q.Logical.Validate(); err != nil {
			return fmt.Errorf("criteria.Query.Logical: %v", err)
		}
	}

	// If the filters are more than one we return an error if the operator is empty for more than one filter must be always specified
	if len(q.Filters) > 1 {
		if q.Logical.String() == "" {
			return fmt.Errorf("criteria.Query.Logical: Logical operator is required for more than 1 filter")
		}
		// If the Logical operator is present we validate if it´s a valid Logical Operator (Check the Logical.Validate() Method)
		if err := q.Logical.Validate(); err != nil {
			return fmt.Errorf("criteria.Query.Logical: %v", err)
		}
	}

	return nil
}
