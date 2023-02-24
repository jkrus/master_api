package files_i

import (
	"github.com/jkrus/master_api/internal/stores/db/repo/files"

	"gorm.io/gorm"
)

type FilesDBStore struct {
	File       files.IFileRepository
	FileStatus files.FileStatusRepositoryI
	FileType   files.FileTypeRepositoryI
}

func NewFilesDBStore(dbHandler *gorm.DB) *FilesDBStore {
	return &FilesDBStore{
		File:       files.NewFileRepository(dbHandler),
		FileStatus: files.NewFileStatusRepository(dbHandler),
		FileType:   files.NewFileTypeRepository(dbHandler),
	}
}
