// Package models provides data models and utility functions for orders.
package models

import (
	"errors"
	"fmt"
	"strings"
)

// Order represents the order type.
type Order string

// Predefined order constants.
const (
	ASCOrder  Order = "asc"
	DESCOrder Order = "desc"
)

var validOrders = map[string]bool{
	ASCOrder.String():  true,
	DESCOrder.String(): true,
}

// NewOrder creates a new Order based on the given string.
// The string is converted to lowercase before creating the Order.
func NewOrder(s string) (Order, error) {
	o := Order(strings.ToLower(s))
	if err := o.Validate(); err != nil {
		return o, err
	}
	return o, nil
}

// Equals checks if the current Order is equal to the provided Order.
// It compares the string representations of the Orders.
func (or Order) Equals(other Order) bool {
	return or.String() == other.String()
}

// String returns the string representation of the Order.
func (or Order) String() string {
	return string(or)
}

// Validate checks if the Order is valid.
// It verifies if the Order exists in the validOrders map.
func (or Order) Validate() error {
	if or.String() == "" {
		return errors.New("invalid order: empty string")
	}

	if !validOrders[or.String()] {
		return fmt.Errorf("invalid order [available:(asc,desc)]: %s", or)
	}
	return nil
}
