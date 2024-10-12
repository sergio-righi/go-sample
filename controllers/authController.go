package controllers

import (
	"encoding/json"
	"net/http"

	"go-sample/models"
	"go-sample/services"
	"go-sample/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthControllerType struct {
	JwtUtils *utils.JwtUtilsType
	Base     *BaseControllerType[models.User, services.UserServiceType]
}

func AuthController(collection *mongo.Collection) *AuthControllerType {
	return &AuthControllerType{
		Base:     BaseController[models.User](collection, services.UserService(collection)),
		JwtUtils: utils.JwtUtils(),
	}
}

func (auth *AuthControllerType) Auth(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Find the user in MongoDB
	var user models.User
	err = auth.Base.Model.FindOne(r.Context(), bson.M{"username": creds.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate a new token
	claims := &models.Claims{}
	tokenString, err := auth.JwtUtils.Generate(w, claims)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Respond with the new token
	auth.Base.HttpControllerType.JsonResponse(w, r, tokenString, http.StatusOK)
}
