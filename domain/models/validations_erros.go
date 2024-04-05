package models

import (
	"strings"
)

// ValidationErrors represents a slice of errors that occurred during validation.
type ValidationErrors []error

// Error returns a string representation of the validation errors.
func (ve ValidationErrors) Error() string {
	var errorMessages []string
	for _, err := range ve {
		errorMessages = append(errorMessages, err.Error())
	}
	return strings.Join(errorMessages, ", ")
}
