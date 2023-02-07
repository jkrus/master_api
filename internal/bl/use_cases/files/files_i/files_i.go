package files_i

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files"
)

type Logic struct {
	File files.FilesI
}

func NewFileLogic(di internal.IAppDeps) *Logic {
	return &Logic{
		File: files.NewFilesI(di),
	}
}
