package bl

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files/files_i"
	"github.com/jkrus/master_api/internal/bl/use_cases/ping/ping_i"
)

type BL struct {
	Ping            *ping_i.Logic
	File            *files_i.Logic
	FileStatusLogic *files_i.FileStatusLogic
}

func NewBL(di internal.IAppDeps) *BL {
	return &BL{
		Ping:            ping_i.NewPingLogic(di),
		File:            files_i.NewFileLogic(di),
		FileStatusLogic: files_i.NewFileStatusLogic(di),
	}
}
