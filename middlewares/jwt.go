package middlewares

import (
	"go-sample/controllers"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type JwtType struct {
	authController *controllers.AuthControllerType
}

func Jwt(authController *controllers.AuthControllerType) *JwtType {
	return &JwtType{authController}
}

func (tk *JwtType) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the request cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		claims := &controllers.Claims{}

		// Parse the JWT token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return tk.authController.JwtSecretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
