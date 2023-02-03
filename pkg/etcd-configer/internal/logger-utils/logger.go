package loggerutils

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewMinimalLoggerConfig ...
//
// It's used for suppressing logs from the etcd client.
//
func NewMinimalLoggerConfig() *zap.Config {
	return &zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.FatalLevel),
		Encoding: "json",
	}
}

// NewTestLogger ...
func NewTestLogger(writer io.Writer) *zap.Logger {
	return zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		}),
		zapcore.AddSync(writer),
		zapcore.DebugLevel,
	))
}
