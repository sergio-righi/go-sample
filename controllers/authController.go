package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go-sample/models"
	"go-sample/services"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthControllerType struct {
	JwtSecretKey string
	Base         *BaseControllerType[models.User, services.UserServiceType]
}

func AuthController(collection *mongo.Collection) *AuthControllerType {
	userService := services.UserService(collection)
	return &AuthControllerType{
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
		Base:         BaseController[models.User](collection, userService),
	}
}

func (auth *AuthControllerType) Login(w http.ResponseWriter, r *http.Request) {
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

	// Create JWT token (same as before)
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(auth.JwtSecretKey))
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
