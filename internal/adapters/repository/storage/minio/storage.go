package minio

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/Gandalf-Le-Dev/quip/internal/core/domain"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
	log        *slog.Logger
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool, log *slog.Logger) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		log.Info("Bucket not found, creating it", "bucket", bucketName)
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
		log:        log,
	}, nil
}

func (s *MinioStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	logger := s.log.With("bucket", s.bucketName, "key", key)
	logger.Debug("Uploading file to Minio")
	_, err := s.client.PutObject(ctx, s.bucketName, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		logger.Error("Failed to upload to Minio", "error", err)
		return err
	}
	logger.Debug("File uploaded successfully to Minio")
	return nil
}

func (s *MinioStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	logger := s.log.With("bucket", s.bucketName, "key", key)
	logger.Debug("Attempting to download file from Minio")
	obj, err := s.client.GetObject(ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		// Check if the error is a Minio ErrorResponse and if it's a NoSuchKey error
		if minioErr, ok := err.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			logger.Warn("Object not found in Minio", "error", err)
			return nil, domain.ErrNotFound
		}
		logger.Error("Failed to get object from Minio", "error", err)
		return nil, err
	}
	logger.Debug("Successfully retrieved object from Minio")
	return obj, nil
}

func (s *MinioStorage) Delete(ctx context.Context, key string) error {
	s.log.Debug("Deleting file from Minio", "bucket", s.bucketName, "key", key)
	return s.client.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{})
}

func (s *MinioStorage) GetURL(ctx context.Context, key string) (string, error) {
	// For public access or generate presigned URL
	return fmt.Sprintf("/files/download/%s", key), nil
}
