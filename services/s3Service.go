package services

import (
	"io"
	"log"
	"time"

	"go-sample/integrations"
)

type S3ServiceType struct {
	MinIO *integrations.MinioType
}

func S3Service(minio *integrations.MinioType) *S3ServiceType {
	return &S3ServiceType{
		MinIO: minio,
	}
}

func (s3 *S3ServiceType) Upload(bucketName, objectName string, file io.Reader, fileSize int64) error {
	err := s3.MinIO.Upload(bucketName, objectName, file, fileSize)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return err
	}
	return nil
}

func (s3 *S3ServiceType) Get(bucketName, objectName string) (io.ReadSeekCloser, error) {
	reader, err := s3.MinIO.Get(bucketName, objectName)
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		return nil, err
	}
	return reader, nil
}

func (s3 *S3ServiceType) List(bucketName, prefix string) ([]string, error) {
	objects, err := s3.MinIO.List(bucketName, prefix, true)
	if err != nil {
		log.Printf("Failed to list files: %v", err)
		return nil, err
	}

	var files []string
	for _, obj := range objects {
		files = append(files, obj.Key)
	}
	return files, nil
}

func (s3 *S3ServiceType) Delete(bucketName, objectName string) error {
	err := s3.MinIO.Delete(bucketName, objectName)
	if err != nil {
		log.Printf("Failed to delete file: %v", err)
		return err
	}
	return nil
}

func (s3 *S3ServiceType) Share(bucketName, objectName string, expiry time.Duration) (string, error) {
	url, err := s3.MinIO.Share(bucketName, objectName, expiry)
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		return "", err
	}
	return url, nil
}

// Rename renames an object in the S3 bucket
func (s3 *S3ServiceType) Rename(bucketName, oldObjectName, newObjectName string) error {
	// Copy the old object to a new one
	err := s3.MinIO.Rename(bucketName, oldObjectName, newObjectName)
	if err != nil {
		return err
	}

	return nil
}
