package models

import (
	"fmt"
)

// Operator is the operator that evaluate the value to the field.
type Operator string

const (
	EqualsOperator       Operator = "="
	NotEqualsOperator    Operator = "!="
	GreaterThan          Operator = ">"
	LessThan             Operator = "<"
	GreaterAndEqualsThan Operator = ">="
	LessAndEqualsThan    Operator = "<="
)

var validOperators = map[string]Operator{
	EqualsOperator.String():       EqualsOperator,
	NotEqualsOperator.String():    NotEqualsOperator,
	GreaterThan.String():          GreaterThan,
	LessThan.String():             LessThan,
	GreaterAndEqualsThan.String(): GreaterAndEqualsThan,
	LessAndEqualsThan.String():    LessAndEqualsThan,
}

func NewOperator(s string) (Operator, error) {
	o := Operator(s)
	if err := o.Validate(); err != nil {
		return o, err
	}
	return o, nil
}

func (o Operator) Equals(other Operator) bool {
	return o.String() == other.String()
}

func (o Operator) String() string {
	return string(o)
}

func (o Operator) Validate() error {
	if o.String() == "" {
		return fmt.Errorf("invalid operator: empty operator")
	}
	_, ok := validOperators[o.String()]
	if !ok {
		return fmt.Errorf("invalid operator: %s", o.String())
	}
	return nil
}
