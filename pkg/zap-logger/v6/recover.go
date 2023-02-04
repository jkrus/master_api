package zaplogger

import (
	"go.uber.org/zap"
)

// Recover ...
func Recover(logger *zap.Logger) {
	if value := recover(); value != nil {
		logger.
			With(zap.Any("panic", value), zap.Stack("stacktrace")).
			Error("panic has been occurred")
	}
}
