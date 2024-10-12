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
	return &UserControllerType{
		Base: BaseController[models.User](collection, services.UserService(collection)),
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
	if err := controller.Base.Service.Create(user); err != nil {
		controller.Base.JsonResponse(w, r, map[string]string{"error": "Could not create user: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	// Respond with success using JsonResponse
	controller.Base.JsonResponse(w, r, user, http.StatusCreated)
}
