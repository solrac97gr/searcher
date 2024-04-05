package models

// SuperFilter are filters that will be applied in the top of your query ignoring validation for permitted search fields it can only be applied as a equals filter
type SuperFilter struct {
	Field string      `bson:"field"`
	Value interface{} `bson:"value"`
}
