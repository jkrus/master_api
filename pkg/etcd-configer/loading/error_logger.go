package loading

import (
	"go.uber.org/zap"
)

// ErrorLogger ...
type ErrorLogger struct {
	baseLogger  *zap.Logger
	baseMessage string
}

// NewErrorLogger ...
func NewErrorLogger(baseLogger *zap.Logger, baseMessage string) ErrorLogger {
	return ErrorLogger{baseLogger, baseMessage}
}

// HandleError ...
func (logger ErrorLogger) HandleError(err error) {
	logger.baseLogger.With(zap.Error(err)).Error(logger.baseMessage)
}
