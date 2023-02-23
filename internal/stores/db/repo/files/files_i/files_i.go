package files_i

import (
	"github.com/jkrus/master_api/internal/stores/db/repo/files"

	"gorm.io/gorm"
)

type FilesDBStore struct {
	FileRepository       files.IFileRepository
	FileStatusRepository files.FileStatusRepositoryI
	FileTypeRepository   files.FileTypeRepositoryI
}

func NewFilesDBStore(dbHandler *gorm.DB) *FilesDBStore {
	return &FilesDBStore{
		FileRepository:       files.NewFileRepository(dbHandler),
		FileStatusRepository: files.NewFileStatusRepository(dbHandler),
		FileTypeRepository:   files.NewFileTypeRepository(dbHandler),
	}
}
