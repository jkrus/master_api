package files_i

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files"
)

type FileTypeLogic struct {
	FileType files.FileTypesI
}

func NewFileTypeLogic(di internal.IAppDeps) *FileTypeLogic {
	return &FileTypeLogic{
		FileType: files.NewFileTypeI(di),
	}
}
