package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoQuery map[string]interface{}

func (mq MongoQuery) GetFilters() (bson.M, error) {
	filters, ok := mq["filters"].(bson.M)
	if !ok {
		return nil, errors.New("invalid mongo filters")
	}
	return filters, nil
}

func (mq MongoQuery) GetSorts() (bson.M, error) {
	filters, ok := mq["sorts"].(bson.M)
	if !ok {
		return nil, errors.New("invalid mongo sorts")
	}
	return filters, nil
}
