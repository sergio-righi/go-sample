package middlewares

import (
	"go-sample/controllers"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type JwtType struct {
	AuthController *controllers.AuthControllerType
}

func Jwt(authController *controllers.AuthControllerType) *JwtType {
	return &JwtType{
		AuthController: authController,
	}
}

func (tk *JwtType) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		// Check for token in Cookie
		cookie, err := r.Cookie("token")
		if err == nil {
			tokenStr = cookie.Value
		} else {
			// If no cookie, check Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			tokenStr = authHeader[len("Bearer "):] // Extract the token part
		}

		claims := &controllers.Claims{}

		// Parse the JWT token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(tk.AuthController.JwtSecretKey), nil // Ensure this is a byte slice
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
