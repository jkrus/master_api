package db

import (
	"database/sql"

	"gorm.io/gorm"

	"github.com/jkrus/master_api/internal/stores/db/repo/files"
)

// DBRepo - интерфейс работы с базой данных
type DBRepo struct {
	DB                   *gorm.DB
	FileRepository       files.IFileRepository
	FileStatusRepository files.FileStatusRepositoryI
}

// NewDBRepo - конструктор интерфейса работы с базой данных
func NewDBRepo(dbHandler *gorm.DB) *DBRepo {
	return &DBRepo{
		DB:                   dbHandler,
		FileRepository:       files.NewFileRepository(dbHandler),
		FileStatusRepository: files.NewFileStatusRepository(dbHandler),
	}
}

// ApplyAutoMigrations - регистрация авто миграции схемы бд из моделей
func ApplyAutoMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&files.File{},
		&files.FileStatus{},
	)
}

// WithTransaction - обертка заворачивающая выполнение операция GORM в транзакцию
func (ds *DBRepo) WithTransaction(handler func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {

	err := ds.DB.Transaction(handler, opts...)
	if err != nil {
		return err
	}
	return nil
}
