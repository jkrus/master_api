package minio

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/jkrus/master_api/internal/config"
)

type Object struct {
	ID   string
	Size int64
	Tags map[string]string
}

func NewClient(config *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(config.MinioEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}

	return minioClient, nil
}
