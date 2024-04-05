package models

import "fmt"

// Filters represents a collection of filters.
type Filters []Filter

// Validate checks the validity of each filter in the collection.
// It returns a ValidationErrors slice if any filters fail validation.
func (fs Filters) Validate() error {
	var validationErrors ValidationErrors
	for index, filter := range fs {
		if err := filter.Validate(); err != nil {
			newErr := fmt.Errorf("filter[%v]: %v", index, err)
			validationErrors = append(validationErrors, newErr)
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func (fs Filters) Len() int { return len(fs) }

// Filter represents a filter this contain a group of conditions operated by a logical operator.
type Filter struct {
	// Conditions is a list of conditions that will conform to this filter
	Conditions Conditions
	// Logical is the logic operation to apply to the group of conditions
	Logical Logical
}

// Validate checks the validity of the filter.
// It returns a ValidationErrors slice if any validation rules fail.
func (f Filter) Validate() error {
	var validationErrors ValidationErrors

	if f.Logical.String() != "" {
		if err := f.Logical.Validate(); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if len(f.Conditions) > 1 {
		if f.Logical.String() == "" {
			return fmt.Errorf("filter.Logical: Logical operator is required for more than 1 condition")
		}
		// If the Logical operator is present we validate if it´s a valid Logical Operator (Check the Logical.Validate() Method)
		if err := f.Logical.Validate(); err != nil {
			return err
		}
	}

	if err := f.Conditions.Validate(); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
