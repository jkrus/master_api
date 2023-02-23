package files_i

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files"
)

type FileStatusLogic struct {
	FileStatus files.FileStatusesI
}

func NewFileStatusLogic(di internal.IAppDeps) *FileStatusLogic {
	return &FileStatusLogic{
		FileStatus: files.NewFileStatusI(di),
	}
}
