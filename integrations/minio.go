package integrations

import (
	"context"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioType struct{ minioClient *minio.Client }

func Minio() *MinioType {
	return &MinioType{}
}

func (s3 *MinioType) Connect() {
	var err error

	// Get MinIO credentials from environment
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	s3.minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to MinIO!")
}

func (s3 *MinioType) Upload(bucketName, objectName string, reader io.Reader, fileSize int64) error {
	// Upload the file
	n, err := s3.minioClient.PutObject(context.Background(), bucketName, objectName, reader, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	log.Printf("Uploaded %s of size %d", objectName, n)
	return nil
}

func (s3 *MinioType) Get(bucketName, objectName string) (io.ReadSeekCloser, error) {
	obj, err := s3.minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// ListBucketFiles lists all objects in the bucket with an optional prefix and delimiter
func (s3 *MinioType) List(bucketName, prefix string, recursive bool) ([]minio.ObjectInfo, error) {
	// Define the context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Prepare the options for listing objects
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	}

	// Get the list of objects
	objectsCh := s3.minioClient.ListObjects(ctx, bucketName, opts)

	// Collect the objects in a slice
	var objects []minio.ObjectInfo
	for object := range objectsCh {
		if object.Err != nil {
			log.Printf("Error listing object: %v", object.Err)
			return nil, object.Err
		}
		objects = append(objects, object)
	}

	return objects, nil
}

// DeleteFile deletes a file from a bucket
func (s3 *MinioType) Delete(bucketName, objectName string) error {
	err := s3.minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	log.Printf("Deleted object %s from bucket %s", objectName, bucketName)
	return nil
}

// ShareFile generates a presigned URL to share a file (valid for a specified duration)
func (s3 *MinioType) Share(bucketName, objectName string, expiry time.Duration) (string, error) {
	// Create an empty url.Values (request parameters)
	reqParams := make(url.Values)

	// Generate a presigned URL
	presignedURL, err := s3.minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}

	log.Printf("Generated presigned URL for object %s in bucket %s", objectName, bucketName)
	return presignedURL.String(), nil
}

// RenameFile renames a file by copying it to a new object and then deleting the old object
func (s3 *MinioType) Rename(bucketName, oldName, newName string) error {
	// Copy object to new location
	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: oldName,
	}
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: newName,
	}

	_, err := s3.minioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		return err
	}

	// Delete the old object
	err = s3.Delete(bucketName, oldName)
	if err != nil {
		return err
	}

	log.Printf("Renamed object from %s to %s in bucket %s", oldName, newName, bucketName)
	return nil
}
