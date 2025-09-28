package image

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/report"
)

type repo struct {
	client     *minio.Client
	bucketName string
	endpoint   string
}

type Repo interface {
	Save(ctx context.Context, ext, contentType string, content io.Reader) (report.ImageURL, error)
	Delete(ctx context.Context, url report.ImageURL) error
}

func NewImageRepoMinio(client *minio.Client, endpoint, bucketName string) Repo {
	return &repo{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
	}
}

func (r *repo) Save(ctx context.Context, ext, contentType string, content io.Reader) (report.ImageURL, error) {
	id := uuid.New().String()
	name := fmt.Sprintf("%s%s", id, ext)

	userMetadata := map[string]string{"x-amz-acl": "public-read"}

	_, err := r.client.PutObject(ctx, r.bucketName, name, content, -1, minio.PutObjectOptions{
		UserMetadata: userMetadata,
		ContentType:  contentType,
	})

	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", r.endpoint, r.bucketName, name)

	return report.ImageURL(url), nil
}

func (r *repo) Delete(ctx context.Context, url report.ImageURL) error {
	parts := strings.Split(string(url), "/")
	name := parts[len(parts)-1]

	err := r.client.RemoveObject(ctx, r.bucketName, name, minio.RemoveObjectOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}
