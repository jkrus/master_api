package internal

import (
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/stores/db"
	"github.com/jkrus/master_api/internal/stores/minio"
)

// IAppDeps - dependency injection container
type IAppDeps interface {
	MinioRepo() *minio.MinioRepo
	DBRepo() *db.DBRepo
	Log() *zap.Logger
}

type di struct {
	logger  *zap.Logger
	minioDB *minio.MinioRepo
	dbRepo  *db.DBRepo
}

func NewDI(logger *zap.Logger, minioDB *minio.MinioRepo, dbRepo *db.DBRepo) IAppDeps {
	return &di{
		logger:  logger,
		minioDB: minioDB,
		dbRepo:  dbRepo,
	}
}

func (d di) Log() *zap.Logger {
	return d.logger
}

func (d di) MinioRepo() *minio.MinioRepo {
	return d.minioDB
}

func (d di) DBRepo() *db.DBRepo {
	return d.dbRepo
}
