package controllers

import (
	"encoding/json"
	"go-sample/models"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// type BaseControllerType[T any] struct {
// 	Model *mongo.Collection
// }

// func BaseController[T any](model *mongo.Collection) *BaseControllerType[T] {
// 	return &BaseControllerType[T]{
// 		Model: model,
// 	}
// }

type BaseControllerType[T any] struct {
	Model *mongo.Collection
	HttpControllerType
}

func BaseController[T any](model *mongo.Collection) *BaseControllerType[T] {
	return &BaseControllerType[T]{
		Model:              model,
		HttpControllerType: *HttpController(),
	}
}

func (bc *BaseControllerType[T]) fromHex(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objectID, nil
}

func (bc *BaseControllerType[T]) jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (bc *BaseControllerType[T]) errorResponse(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
}

func (bc *BaseControllerType[T]) All(w http.ResponseWriter, r *http.Request) {
	cursor, err := bc.Model.Find(r.Context(), bson.D{{}})
	if err != nil {
		bc.errorResponse(w, "Failed to find documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var results []models.User // Change to your specific model type
	if err := cursor.All(r.Context(), &results); err != nil {
		bc.errorResponse(w, "Failed to parse documents", http.StatusInternalServerError)
		return
	}

	bc.jsonResponse(w, results, http.StatusOK)
}

func (bc *BaseControllerType[T]) Create(w http.ResponseWriter, r *http.Request) {
	var doc T // Change to your specific model type
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		bc.errorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := bc.Model.InsertOne(r.Context(), doc)
	if err != nil {
		bc.errorResponse(w, "Failed to create document", http.StatusInternalServerError)
		return
	}

	bc.jsonResponse(w, doc, http.StatusCreated)
}

func (bc *BaseControllerType[T]) Find(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"] // Extract the ID from the URL

	objectID, err := bc.fromHex(id)
	if err != nil {
		bc.errorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var result T // Change to your specific model type
	err = bc.Model.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		bc.errorResponse(w, "Document not found", http.StatusNotFound)
		return
	}

	bc.jsonResponse(w, result, http.StatusOK)
}

func (bc *BaseControllerType[T]) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"] // Extract the ID from the URL

	objectID, err := bc.fromHex(id)
	if err != nil {
		bc.errorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var doc T // Change to your specific model type
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		bc.errorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update the document
	_, err = bc.Model.UpdateOne(r.Context(), bson.M{"_id": objectID}, bson.M{"$set": doc})
	if err != nil {
		bc.errorResponse(w, "Failed to update document", http.StatusInternalServerError)
		return
	}

	bc.jsonResponse(w, doc, http.StatusOK)
}

func (bc *BaseControllerType[T]) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"] // Extract the ID from the URL

	objectID, err := bc.fromHex(id)
	if err != nil {
		bc.errorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	_, err = bc.Model.DeleteOne(r.Context(), bson.M{"_id": objectID})
	if err != nil {
		bc.errorResponse(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
