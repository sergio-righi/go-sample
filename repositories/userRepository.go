package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryType struct {
	Collection *mongo.Collection
}

func UserRepository(collection *mongo.Collection) *UserRepositoryType {
	return &UserRepositoryType{
		Collection: collection,
	}
}
