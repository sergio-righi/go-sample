package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type PersonRepositoryType struct {
	Collection *mongo.Collection
}

func PersonRepository(collection *mongo.Collection) *PersonRepositoryType {
	return &PersonRepositoryType{
		Collection: collection,
	}
}
