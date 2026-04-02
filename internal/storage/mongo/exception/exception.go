package exception

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// IsDuplicateKey checks if the error is a MongoDB duplicate key error.
func IsDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(mongo.WriteException); ok {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 || we.Code == 11001 {
				return true
			}
		}
	}
	if e, ok := err.(mongo.CommandError); ok {
		return e.Code == 11000 || e.Code == 11001
	}
	// For bulk write exceptions
	if e, ok := err.(mongo.BulkWriteException); ok {
		if e.WriteErrors != nil {
			for _, we := range e.WriteErrors {
				if we.Code == 11000 || we.Code == 11001 {
					return true
				}
			}
		}
	}
	return false
}
