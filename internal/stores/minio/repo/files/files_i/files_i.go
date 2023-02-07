package files_i

import (
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/stores/minio/repo/files"
)

type FileStoreI struct {
	FileStore files.FileStoreI
}

func NewFileStore(logger *zap.Logger, minioClient *minio.Client) FileStoreI {
	return FileStoreI{
		FileStore: files.NewFileStoreI(logger, minioClient),
	}
}
