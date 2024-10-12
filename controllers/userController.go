package controllers

import (
	"encoding/json"
	"go-sample/models"
	"go-sample/services"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserControllerType struct {
	Base *BaseControllerType[models.User, services.UserServiceType]
}

func UserController(collection *mongo.Collection) *UserControllerType {
	userService := services.UserService(collection)
	return &UserControllerType{
		Base: BaseController[models.User](collection, userService),
	}
}

func (controller *UserControllerType) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode the request body into the user model
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword) // Replace plain password with hashed password

	// Call the user service to create the user
	if err := controller.Base.Service.CreateUser(user); err != nil {
		http.Error(w, "Could not create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
