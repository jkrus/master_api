package files_i

import (
	"github.com/jkrus/master_api/internal/stores/db/repo/files"

	"gorm.io/gorm"
)

type FilesDBStore struct {
	FileRepository files.IFileRepository
}

func NewNotificationDBStore(dbHandler *gorm.DB, dbReader *gorm.DB) *FilesDBStore {
	return &FilesDBStore{
		FileRepository: files.NewFileRepository(dbHandler),
	}
}
