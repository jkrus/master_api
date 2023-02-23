package files_i

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files"
)

type Logic struct {
	File       files.FilesI
	FileStatus files.FileStatusesI
	FileType   files.FileTypesI
}

func NewFileLogic(di internal.IAppDeps) *Logic {
	return &Logic{
		File:       files.NewFilesI(di),
		FileStatus: files.NewFileStatusI(di),
		FileType:   files.NewFileTypeI(di),
	}
}
