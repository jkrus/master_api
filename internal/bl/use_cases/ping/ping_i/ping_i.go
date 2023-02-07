package ping_i

import (
	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl/use_cases/ping"
)

type Logic struct {
	Ping ping.PingI
}

func NewPingLogic(di internal.IAppDeps) *Logic {
	return &Logic{
		Ping: ping.NewPing(di),
	}
}
