package ping

import (
	"context"

	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/pkg/errors"
)

type PingI interface {
	Ping(ctx context.Context, msg string) (string, error)
}

type ping struct {
	di internal.IAppDeps
}

func NewPing(di internal.IAppDeps) PingI {
	return &ping{di: di}
}

func (p ping) Ping(ctx context.Context, msg string) (string, error) {
	if msg == "ping" {
		return "pong", nil
	}
	return "", errors.Ctx().Just(errors.New("wrong ping request"))
}
