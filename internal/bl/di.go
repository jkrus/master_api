package bl

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/files/files_i"
	"github.com/jkrus/master_api/internal/bl/use_cases/ping/ping_i"
	"github.com/jkrus/master_api/internal/bl/use_cases/users/users_i"
)

type BL struct {
	Ping       *ping_i.Logic
	FileLogic  *files_i.Logic
	UsersLogic *users_i.UserLogic
}

func NewBL(di internal.IAppDeps) *BL {
	return &BL{
		Ping:       ping_i.NewPingLogic(di),
		FileLogic:  files_i.NewFileLogic(di),
		UsersLogic: users_i.NewUserLogic(di),
	}
}
