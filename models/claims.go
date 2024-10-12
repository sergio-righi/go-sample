package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
