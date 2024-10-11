package controllers

import (
	"encoding/json"
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
func (hc *HttpControllerType) JsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// errorResponse to send error messages in the response
func (hc *HttpControllerType) ErrorResponse(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
}
