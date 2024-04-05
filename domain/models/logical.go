package models

import (
	"fmt"
	"strings"
)

// Logical is the logic operation to apply to the group of filters
type Logical string

const (
	ANDLogical Logical = "and"
	ORLogical  Logical = "or"
)

var validLogical = map[string]Logical{
	ANDLogical.String(): ANDLogical,
	ORLogical.String():  ORLogical,
}

// NewLogical creates a new Logical based on the given string.
// If the string is not a valid logical operator, it returns an error.
func NewLogical(s string) (Logical, error) {
	//To lowercase the operator
	s = strings.ToLower(s)
	l := Logical(s)
	if err := l.Validate(); err != nil {
		return "", err
	}
	return l, nil
}

// Equals checks if the current Logical is equal to the provided Logical.
func (l Logical) Equals(other Logical) bool {
	return l == other
}

// String returns the string representation of the Logical.
func (l Logical) String() string {
	return string(l)
}

// Validate checks the validity of the Logical.
func (l Logical) Validate() error {
	_, ok := validLogical[l.String()]
	if !ok {
		if l.String() == "" {
			return fmt.Errorf("invalid Logical operator cannot be empty")
		}
		return fmt.Errorf("invalid logical operator: %s", l.String())
	}
	return nil
}
