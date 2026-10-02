package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMinioUploadReadAndRemove(t *testing.T) {
	endpoint := os.Getenv("TEST_MINIO_ENDPOINT")
	if endpoint == "" {
		t.Skip("TEST_MINIO_ENDPOINT is not set")
	}
	t.Setenv("MINIO_ENDPOINT", endpoint)
	s, err := NewMinio()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	key := "integration-test/" + uuid.NewString()
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Remove(cleanupCtx, key); err != nil {
			t.Errorf("cleanup: %v", err)
		}
	})
	want := []byte("media storage integration test")
	if err := s.Upload(ctx, key, bytes.NewReader(want), int64(len(want)), "text/plain"); err != nil {
		t.Fatal(err)
	}
	reader, contentType, err := s.Open(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	got, err := io.ReadAll(reader)
	_ = reader.Close()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) || contentType != "text/plain" {
		t.Fatalf("download mismatch: %q, %q", got, contentType)
	}
	if err := s.Remove(ctx, key); err != nil {
		t.Fatal(err)
	}
	_, _, err = s.Open(ctx, key)
	if !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("expected missing object, got %v", err)
	}
}
