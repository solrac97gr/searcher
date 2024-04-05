package models

import (
	"errors"
	"fmt"
	"reflect"
)

// Conditions represents a collection of conditions.
type Conditions []Condition

// Validate checks the validity of each condition in the collection.
// It returns a ValidationErrors slice if any conditions fail validation.
func (cs Conditions) Validate() error {
	if len(cs) == 0 {
		return errors.New("empty Conditions: at least one condition must be specified")
	}

	var validationErrors ValidationErrors
	for index, condition := range cs {
		if err := condition.Validate(); err != nil {
			newErr := fmt.Errorf("condition[%v]: %v", index, err)
			validationErrors = append(validationErrors, newErr)
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

// Len returns the number of elements in the collection
func (cs Conditions) Len() int { return len(cs) }

// Condition represents a single condition.
type Condition struct {
	// Field is the name of the field you wanna evaluate the condition against.
	Field Field `json:"field" example:"amount"`
	// Operator is the operator that evaluate the value to the field.
	Operator Operator `json:"operator" example:">="`
	// Value is the value to evaluate the field.
	Value interface{} `json:"value"`
}
type emptyStruct struct{}

// Validate checks the validity of the condition.
// It returns a ValidationErrors slice if any validation rules fail.
func (c Condition) Validate() error {
	var validationErrors ValidationErrors
	if err := c.Field.Validate(); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if c.Value == nil {
		validationErrors = append(validationErrors, errors.New("invalid value: cannot be nil"))
	} else {
		valueType := reflect.TypeOf(c.Value)
		if valueType.Kind() == reflect.Struct && valueType == reflect.TypeOf(emptyStruct{}) {
			validationErrors = append(validationErrors, errors.New("invalid value: cannot be an empty struct"))
		}

		if valueType.Kind() == reflect.Map {
			if len(c.Value.(map[string]interface{})) == 0 {
				validationErrors = append(validationErrors, errors.New("invalid value: cannot be empty map"))
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
