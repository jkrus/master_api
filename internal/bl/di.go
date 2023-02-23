package bl

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files/files_i"
	"github.com/jkrus/master_api/internal/bl/use_cases/ping/ping_i"
)

type BL struct {
	Ping      *ping_i.Logic
	FileLogic *files_i.Logic
}

func NewBL(di internal.IAppDeps) *BL {
	return &BL{
		Ping:      ping_i.NewPingLogic(di),
		FileLogic: files_i.NewFileLogic(di),
	}
}
