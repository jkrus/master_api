package minio

import (
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/stores/minio/repo/files/files_i"
)

// MinioRepo - интерфейс работы с MinIO
type MinioRepo struct {
	Files files_i.FileStoreI
}

// NewMinioRepo - конструктор интерфейса работы с MinIO
func NewMinioRepo(logger *zap.Logger, minioClient *minio.Client) *MinioRepo {
	return &MinioRepo{
		Files: files_i.NewFileStore(logger, minioClient),
	}
}
