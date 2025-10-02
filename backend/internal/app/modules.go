package app

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/config"
)

func initPostgresClient(config *config.PostgresConfig) *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database)
	postgresClient, err := sqlx.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	return postgresClient
}

func initMinioConnection(cfg *config.MinioConfig) (*minio.Client, error) {
	ctx := context.Background()

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.RootUser, cfg.RootPassword, ""),
		Secure: false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to minio: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to minio: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create minio bucket: %w", err)
		}
	}

	policy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": "*",
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
            }
        ]
    }`

	if err := client.SetBucketPolicy(ctx, cfg.BucketName, fmt.Sprintf(policy, cfg.BucketName)); err != nil {
		return nil, fmt.Errorf("failed to set up minio bucket: %w", err)
	}

	return client, nil
}
