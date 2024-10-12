package controllers

import (
	"go-sample/models"
	"go-sample/services"

	"go.mongodb.org/mongo-driver/mongo"
)

type PersonControllerType struct {
	Base *BaseControllerType[models.Person, services.PersonServiceType]
}

func PersonController(collection *mongo.Collection) *PersonControllerType {

	personService := services.PersonService(collection)
	return &PersonControllerType{
		Base: BaseController[models.Person](collection, personService),
	}
}
