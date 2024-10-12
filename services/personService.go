package services

import (
	"go-sample/repositories"

	"go.mongodb.org/mongo-driver/mongo"
)

type PersonServiceType struct {
	Repo *repositories.PersonRepositoryType
}

func PersonService(collection *mongo.Collection) *PersonServiceType {
	return &PersonServiceType{
		Repo: repositories.PersonRepository(collection),
	}
}
