package utils

import (
	"go-sample/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtUtilsType struct {
	JwtSecretKey string
}

func JwtUtils() *JwtUtilsType {
	return &JwtUtilsType{
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}

func (utils *JwtUtilsType) Generate(w http.ResponseWriter, claims *models.Claims) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString([]byte(utils.JwtSecretKey))
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return tokenString, nil
}
