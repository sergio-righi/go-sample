package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go-sample/models"

	"github.com/dgrijalva/jwt-go"
)

// Struct to hold the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthControllerType struct {
	JwtSecretKey string
}

func AuthController() *AuthControllerType {
	return &AuthControllerType{
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}

// Dummy login handler for generating JWT token
func (auth *AuthControllerType) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// In a real application, you'd authenticate against your database.
	if creds.Username != "user" || creds.Password != "password" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create JWT token with expiration time of 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(auth.JwtSecretKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// Return the JWT token
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
