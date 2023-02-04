package i_http

import "go.uber.org/zap"

type IHTTPServer interface {
	Run(logger *zap.Logger)
}
