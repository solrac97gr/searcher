package models

import "errors"

// Field represents a field name.
type Field string

// NewField creates a new Field based on the given string.
func NewField(s string) (Field, error) {
	f := Field(s)
	if err := f.Validate(); err != nil {
		return f, err
	}
	return f, nil
}

// FieldType represents the type of a field.
type FieldType string

// Predefined field types.
const (
	Undefined FieldType = "undefined"
	String    FieldType = "string"
	Number    FieldType = "number"
	Date      FieldType = "date"
)

func (ft FieldType) String() string {
	return string(ft)
}

func (ft FieldType) Equals(other FieldType) bool {
	return ft.String() == other.String()
}

// FieldMetaData represents metadata for a field.
type FieldMetaData struct {
	Field      Field
	Type       FieldType
	IsAnalyzed bool
}

// Validate checks the validity of the Field.
func (f Field) Validate() error {
	if f.String() == "" {
		return errors.New("invalid field: cannot be empty")
	}
	return nil
}

// String returns the string representation of the Field.
func (f Field) String() string {
	return string(f)
}

// Equals checks if the current Field is equal to the provided Field.
func (f Field) Equals(other Field) bool {
	return f.String() == other.String()
}
