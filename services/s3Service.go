package services

import (
	"io"
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
	return s3.MinIO.Upload(bucketName, objectName, file, fileSize)
}

func (s3 *S3ServiceType) Get(bucketName, objectName string) (io.ReadSeekCloser, error) {
	reader, err := s3.MinIO.Get(bucketName, objectName)
	return reader, err
}

func (s3 *S3ServiceType) List(bucketName, prefix string, recursive bool) ([]string, error) {
	objects, err := s3.MinIO.List(bucketName, prefix, recursive)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range objects {
		files = append(files, obj.Key)
	}
	return files, nil
}

func (s3 *S3ServiceType) Delete(bucketName, objectName string) error {
	return s3.MinIO.Delete(bucketName, objectName)
}

func (s3 *S3ServiceType) Share(bucketName, objectName string, expiry time.Duration) (string, error) {
	return s3.MinIO.Share(bucketName, objectName, expiry)
}

func (s3 *S3ServiceType) Rename(bucketName, oldObjectName, newObjectName string) error {
	return s3.MinIO.Rename(bucketName, oldObjectName, newObjectName)
}
