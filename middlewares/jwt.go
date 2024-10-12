package middlewares

import (
	"context"
	"go-sample/controllers"
	"go-sample/models"
	"go-sample/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type JwtType struct {
	AuthController *controllers.AuthControllerType
	JwtUtils       *utils.JwtUtilsType
}

func Jwt(authController *controllers.AuthControllerType) *JwtType {
	return &JwtType{
		AuthController: authController,
		JwtUtils:       utils.JwtUtils(),
	}
}

// Middleware that validates the JWT token and refreshes it if needed
func (tk *JwtType) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		// Check for token in Cookie or Authorization header
		cookie, err := r.Cookie("token")
		if err == nil {
			tokenStr = cookie.Value
		} else {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}
			tokenStr = authHeader[len("Bearer "):] // Extract the token part
		}

		// Parse the JWT token
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(tk.JwtUtils.JwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if the token is close to expiration (within 1 minute)
		var newTokenString string
		// Refresh the token if near expiration
		newTokenString, err = tk.JwtUtils.Generate(w, claims)
		if err != nil {
			http.Error(w, "Could not refresh token", http.StatusInternalServerError)
			return
		}

		// Store the claims and the new token in context
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		if newTokenString != "" {
			ctx = context.WithValue(ctx, "refreshedToken", newTokenString)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
