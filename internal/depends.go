package internal

import (
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/stores/db"
	"github.com/jkrus/master_api/internal/stores/hyper_ledger"
	"github.com/jkrus/master_api/internal/stores/minio"
)

// IAppDeps - dependency injection container
type IAppDeps interface {
	MinioRepo() *minio.MinioRepo
	DBRepo() *db.DBRepo
	HyperLagerStore() *hyper_ledger.HFRepo
	Log() *zap.Logger
}

type di struct {
	logger          *zap.Logger
	minioDB         *minio.MinioRepo
	dbRepo          *db.DBRepo
	hyperLagerStore *hyper_ledger.HFRepo
}

func NewDI(logger *zap.Logger, minioDB *minio.MinioRepo, dbRepo *db.DBRepo, hyperLagerStore *hyper_ledger.HFRepo) IAppDeps {
	return &di{
		logger:          logger,
		minioDB:         minioDB,
		dbRepo:          dbRepo,
		hyperLagerStore: hyperLagerStore,
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

func (d di) HyperLagerStore() *hyper_ledger.HFRepo {
	return d.hyperLagerStore
}
