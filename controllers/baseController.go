package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseControllerType[T any, S any] struct {
	Model   *mongo.Collection
	Service *S
	HttpControllerType
}

func BaseController[T any, S any](model *mongo.Collection, service *S) *BaseControllerType[T, S] {
	return &BaseControllerType[T, S]{
		Model:              model,
		Service:            service,
		HttpControllerType: *HttpController(),
	}
}

func (bc *BaseControllerType[T, S]) All(w http.ResponseWriter, r *http.Request) {
	cursor, err := bc.Model.Find(r.Context(), bson.D{{}})
	if err != nil {
		bc.ErrorResponse(w, "Failed to find documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var results []T
	if err := cursor.All(r.Context(), &results); err != nil {
		bc.ErrorResponse(w, "Failed to parse documents", http.StatusInternalServerError)
		return
	}

	bc.JsonResponse(w, r, results, http.StatusOK)
}

func (bc *BaseControllerType[T, S]) Create(w http.ResponseWriter, r *http.Request) {
	var doc T
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		bc.ErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := bc.Model.InsertOne(r.Context(), doc)
	if err != nil {
		bc.ErrorResponse(w, "Failed to create document", http.StatusInternalServerError)
		return
	}

	bc.JsonResponse(w, r, doc, http.StatusCreated)
}

func (bc *BaseControllerType[T, S]) Find(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"] // Extract the ID from the URL

	objectID, err := bc.FromHex(id)
	if err != nil {
		bc.ErrorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var result T
	err = bc.Model.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		bc.ErrorResponse(w, "Document not found", http.StatusNotFound)
		return
	}

	bc.JsonResponse(w, r, result, http.StatusOK)
}

func (bc *BaseControllerType[T, S]) Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"] // Extract the ID from the URL

	objectID, err := bc.FromHex(id)
	if err != nil {
		bc.ErrorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var doc T
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		bc.ErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update the document
	_, err = bc.Model.UpdateOne(r.Context(), bson.M{"_id": objectID}, bson.M{"$set": doc})
	if err != nil {
		bc.ErrorResponse(w, "Failed to update document", http.StatusInternalServerError)
		return
	}

	bc.JsonResponse(w, r, doc, http.StatusOK)
}

func (bc *BaseControllerType[T, S]) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"] // Extract the ID from the URL

	objectID, err := bc.FromHex(id)
	if err != nil {
		bc.ErrorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	_, err = bc.Model.DeleteOne(r.Context(), bson.M{"_id": objectID})
	if err != nil {
		bc.ErrorResponse(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
