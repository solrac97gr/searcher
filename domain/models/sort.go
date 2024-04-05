package models

import (
	"errors"
	"fmt"
)

type Sorts []Sort

func (ss Sorts) Validate() error {
	for index, sort := range ss {
		if err := sort.Validate(); err != nil {
			newErr := fmt.Errorf("sorts[%v]: %v", index, err)
			return newErr
		}
	}
	return nil
}

type Sort struct {
	Field string
	Order Order
}

func (s Sort) Validate() error {
	if s.Field == "" {
		return errors.New("invalid field: empty")
	}
	if err := s.Order.Validate(); err != nil {
		return err
	}
	return nil
}
