package models

import (
	"errors"
	"fmt"
)

// Pagination represents the structure that indicates pagination to the database
type Pagination struct {
	// Limit is the number of items to be returned
	// - If the limit is 0 or less the default limit value (50) is applied
	Limit uint `json:"limit" example:"100" maximum:"1000" minimum:"1"`
	// Offset is the number of items to be skipped
	Offset uint `json:"offset" example:"100"`
}

const MaximumLimit = 1000
const MaximumLimitOffsetSize uint = 10000

func (p Pagination) Validate() error {
	// The maximum number of items you can retrieve from the database
	if p.Limit > MaximumLimit {
		return errors.New("limit must be less than " + fmt.Sprint(MaximumLimit))
	}
	// This condition is necessary for avoid memory limitations of the database
	if (p.Limit + p.Offset) > MaximumLimitOffsetSize {
		return fmt.Errorf(" limit(%d) + offset(%d) must be less or equals %d", p.Limit, p.Offset, MaximumLimitOffsetSize)
	}
	return nil
}
