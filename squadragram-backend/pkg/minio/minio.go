package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ErrObjectNotFound = errors.New("media object not found")

type MinioStorage struct {
	Client *minio.Client
	Bucket string
}

func NewMinio() (*MinioStorage, error) {
	endpoint := envOrDefault("MINIO_ENDPOINT", "localhost:9005")
	accessKey := envOrDefault("MINIO_ACCESS_KEY", "admin")
	secretKey := envOrDefault("MINIO_SECRET_KEY", "12345678")
	bucket := envOrDefault("MINIO_BUCKET", "squadragram-media")
	secure, err := strconv.ParseBool(envOrDefault("MINIO_SECURE", "false"))
	if err != nil {
		return nil, fmt.Errorf("parse MINIO_SECURE: %w", err)
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			accessKey,
			secretKey,
			"",
		),
		Secure: secure,
	})

	if err != nil {
		return nil, fmt.Errorf("configure MinIO: %w", err)
	}

	return &MinioStorage{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (s *MinioStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	exists, err := s.Client.BucketExists(ctx, s.Bucket)
	if err != nil {
		return fmt.Errorf("check media bucket: %w", err)
	}
	if !exists {
		if err := s.Client.MakeBucket(ctx, s.Bucket, minio.MakeBucketOptions{}); err != nil {
			code := minio.ToErrorResponse(err).Code
			if code != "BucketAlreadyOwnedByYou" {
				return fmt.Errorf("create media bucket: %w", err)
			}
		}
	}
	_, err = s.Client.PutObject(ctx, s.Bucket, key, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("upload media: %w", err)
	}
	return nil
}

func (s *MinioStorage) Remove(ctx context.Context, key string) error {
	return s.Client.RemoveObject(ctx, s.Bucket, key, minio.RemoveObjectOptions{})
}

func (s *MinioStorage) Open(ctx context.Context, key string) (io.ReadCloser, string, error) {
	object, err := s.Client.GetObject(ctx, s.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("get media: %w", err)
	}
	info, err := object.Stat()
	if err != nil {
		_ = object.Close()
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, "", ErrObjectNotFound
		}
		return nil, "", fmt.Errorf("stat media: %w", err)
	}
	return object, info.ContentType, nil
}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
