package internal

import (
	"go.uber.org/zap"
)

// IAppDeps - dependency injection container
type IAppDeps interface{}

type di struct {
	logger *zap.Logger
}

func NewDI(logger *zap.Logger) IAppDeps {
	return &di{
		logger: logger,
	}
}
