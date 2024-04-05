package models

import (
	"fmt"

	"github.com/solrac97gr/searcher/internal/sentinels"
)

// Criteria represents the criteria for a search operation.
// @Description Criteria is the common structure for realize search queries in the endpoints of Sigma Management API.
type Criteria struct {
	// Pagination is the structure that compile the pagination info for the endpoint
	Pagination Pagination `json:"pagination"`
	// Query is the structure that contains the filters and sorts to apply to the search.
	Query Query `json:"query"`
}

// Validate checks the validity of the Criteria.
// It returns a standard validation error (/pkg/errors) if any validation rules fail.
func (c Criteria) Validate() error {
	if err := c.Pagination.Validate(); err != nil {
		return fmt.Errorf("%w: %v", sentinels.ErrValidation, err)
	}
	if err := c.Query.Validate(); err != nil {
		return fmt.Errorf("%w: %v", sentinels.ErrValidation, err)
	}
	return nil
}
