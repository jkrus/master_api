package internal

import (
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/stores/minio"
)

// IAppDeps - dependency injection container
type IAppDeps interface {
	MinioRepo() *minio.MinioRepo
	Log() *zap.Logger
}

type di struct {
	logger  *zap.Logger
	minioDB *minio.MinioRepo
}

func NewDI(logger *zap.Logger, minioDB *minio.MinioRepo) IAppDeps {
	return &di{
		logger:  logger,
		minioDB: minioDB,
	}
}

func (d di) Log() *zap.Logger {
	return d.logger
}

func (d di) MinioRepo() *minio.MinioRepo {
	return d.minioDB
}
