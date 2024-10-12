package controllers

import (
	"go-sample/integrations"
	"go-sample/services"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type S3ControllerType struct {
	Service *services.S3ServiceType
	HttpControllerType
}

func S3Controller(minio *integrations.MinioType) *S3ControllerType {
	return &S3ControllerType{
		Service:            services.S3Service(minio),
		HttpControllerType: *HttpController(),
	}
}

// UploadFile handles the file upload request
func (s3 *S3ControllerType) Upload(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucket")
	objectName := r.URL.Query().Get("object")

	file, _, err := r.FormFile("file")
	if err != nil {
		s3.ErrorResponse(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = s3.Service.Upload(bucketName, objectName, file, r.ContentLength)
	if err != nil {
		s3.ErrorResponse(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	s3.JsonResponse(w, map[string]string{"message": "File uploaded successfully"}, http.StatusCreated)
}

// GetFile handles the file retrieval request
func (s3 *S3ControllerType) Get(w http.ResponseWriter, r *http.Request) {
	bucketName := mux.Vars(r)["bucket"]
	objectName := mux.Vars(r)["object"]

	reader, err := s3.Service.Get(bucketName, objectName)
	if err != nil {
		s3.ErrorResponse(w, "Failed to get file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+objectName)
	http.ServeContent(w, r, objectName, time.Now(), reader)
}

// ListFiles handles the request to list files in a bucket
func (s3 *S3ControllerType) List(w http.ResponseWriter, r *http.Request) {
	bucketName := mux.Vars(r)["bucket"]
	prefix := r.URL.Query().Get("prefix")
	recursive := r.URL.Query().Get("recursive") == "true"

	files, err := s3.Service.List(bucketName, prefix, recursive)
	if err != nil {
		s3.ErrorResponse(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	s3.JsonResponse(w, files, http.StatusOK)
}

// DeleteFile handles the file deletion request
func (s3 *S3ControllerType) Delete(w http.ResponseWriter, r *http.Request) {
	bucketName := mux.Vars(r)["bucket"]
	objectName := mux.Vars(r)["object"]

	err := s3.Service.Delete(bucketName, objectName)
	if err != nil {
		s3.ErrorResponse(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	s3.JsonResponse(w, map[string]string{"message": "File deleted successfully"}, http.StatusNoContent)
}

// ShareFile handles the request to generate a presigned URL for sharing a file
func (s3 *S3ControllerType) Share(w http.ResponseWriter, r *http.Request) {
	bucketName := mux.Vars(r)["bucket"]
	objectName := mux.Vars(r)["object"]
	expiry := time.Duration(24) * time.Hour // Set expiry time for the presigned URL

	url, err := s3.Service.Share(bucketName, objectName, expiry)
	if err != nil {
		s3.ErrorResponse(w, "Failed to generate presigned URL", http.StatusInternalServerError)
		return
	}

	s3.JsonResponse(w, map[string]string{"url": url}, http.StatusOK)
}

// RenameFile handles renaming an object in S3
func (s3 *S3ControllerType) Rename(w http.ResponseWriter, r *http.Request) {
	bucketName := mux.Vars(r)["bucket"]
	oldObjectName := mux.Vars(r)["object"]
	newObjectName := r.URL.Query().Get("new_name")

	err := s3.Service.Rename(bucketName, oldObjectName, newObjectName)
	if err != nil {
		s3.ErrorResponse(w, "Failed to rename file", http.StatusInternalServerError)
		return
	}

	s3.JsonResponse(w, map[string]string{"message": "File renamed successfully"}, http.StatusOK)
}
