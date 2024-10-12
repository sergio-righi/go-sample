package services

import (
	"context"
	"go-sample/models"
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

func (service *UserServiceType) Create(user models.User) error {
	_, err := service.Repo.Collection.InsertOne(context.Background(), user)
	return err
}
