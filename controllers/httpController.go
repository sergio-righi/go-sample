package controllers

import (
	"encoding/json"
	"fmt"
	"go-sample/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HttpControllerType struct{}

func HttpController() *HttpControllerType {
	return &HttpControllerType{}
}

// fromHex method to convert string to MongoDB ObjectID
func (hc *HttpControllerType) FromHex(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objectID, nil
}

// jsonResponse to send JSON responses
func (hc *HttpControllerType) JsonResponse(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	// Check if a refreshed token exists in the request context
	refreshedToken := r.Context().Value("refreshedToken")
	if refreshedToken != nil {
		// Add the refreshed token to the response header
		w.Header().Set("X-Refreshed-Token", refreshedToken.(string))
	}

	// Write the status and response body
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// errorResponse to send error messages in the response
func (hc *HttpControllerType) ErrorResponse(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
}

func JwtHandler(w http.ResponseWriter, r *http.Request) (*models.Claims, error) {
	claims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("unable to extract user information")
	}
	return claims, nil
}
