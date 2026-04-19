package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/digikeys/backend/config"
)

// MinioStorage implements ports.StorageService using MinIO/S3-compatible storage.
type MinioStorage struct {
	client *minio.Client
}

// NewMinioStorage creates a new MinIO storage client and ensures the default bucket exists.
func NewMinioStorage(ctx context.Context, cfg config.StorageConfig) (*MinioStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	// Ensure default bucket exists.
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("check bucket %s: %w", cfg.Bucket, err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("create bucket %s: %w", cfg.Bucket, err)
		}
		slog.Info("created storage bucket", "bucket", cfg.Bucket)
	}

	return &MinioStorage{client: client}, nil
}

// Upload stores an object in the specified bucket and returns the object URL.
func (s *MinioStorage) Upload(ctx context.Context, bucket, key string, reader io.Reader, contentType string) (string, error) {
	if err := s.ensureBucket(ctx, bucket); err != nil {
		return "", err
	}

	_, err := s.client.PutObject(ctx, bucket, key, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload object %s/%s: %w", bucket, key, err)
	}

	return fmt.Sprintf("%s/%s/%s", s.client.EndpointURL(), bucket, key), nil
}

// Download retrieves an object from the specified bucket.
func (s *MinioStorage) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("download object %s/%s: %w", bucket, key, err)
	}
	return obj, nil
}

// Delete removes an object from the specified bucket.
func (s *MinioStorage) Delete(ctx context.Context, bucket, key string) error {
	err := s.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete object %s/%s: %w", bucket, key, err)
	}
	return nil
}

// GetURL returns a presigned URL for the specified object, valid for 1 hour.
func (s *MinioStorage) GetURL(ctx context.Context, bucket, key string) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, bucket, key, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("presign URL %s/%s: %w", bucket, key, err)
	}
	return url.String(), nil
}

// ensureBucket creates the bucket if it does not exist.
func (s *MinioStorage) ensureBucket(ctx context.Context, bucket string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("check bucket %s: %w", bucket, err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("create bucket %s: %w", bucket, err)
		}
	}
	return nil
}
