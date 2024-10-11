package services

import (
	"go-sample/repositories"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceType struct {
	Repo *repositories.UserRepositoryType
}

func UserService(collection *mongo.Collection) *UserServiceType {
	return &UserServiceType{
		Repo: repositories.UserRepository(collection),
	}
}
