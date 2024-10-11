package controllers

import (
	"go-sample/models"
	"go-sample/services"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserControllerType struct {
	Base    *BaseControllerType[models.User]
	Service *services.UserServiceType
}

func UserController(collection *mongo.Collection) *UserControllerType {
	return &UserControllerType{
		Base:    BaseController[models.User](collection),
		Service: services.UserService(collection),
	}
}
