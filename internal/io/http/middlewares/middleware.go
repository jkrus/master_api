package middlewares

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/bl"
)

type IMiddleWares interface {
	GetObservabilityMiddleware(logger *zap.Logger) mux.MiddlewareFunc
}

type middlewares struct {
	bl *bl.BL
}

func NewMiddleWares(bl *bl.BL) IMiddleWares {
	return &middlewares{
		bl: bl,
	}
}
