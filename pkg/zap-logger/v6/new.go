package zaplogger

import (
	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New ...
//
// The app will add as the app field to the log.
//
// The level can be: debug, info, warn, error, dpanic, panic or fatal.
// See for details: https://godoc.org/go.uber.org/zap/zapcore#Level
//
// The encoding can be: console or json.
// See for details: https://godoc.org/go.uber.org/zap/zapcore#Encoder
//
// A case of the level and the encoding isn't important.
//
// An empty value of the level and the encoding is wrong.
//
func New(app, level, encoding string) (*zap.Logger, zap.AtomicLevel, error) {
	parsedLevel, err := ParseLevel(level)
	if err != nil {
		return nil, zap.AtomicLevel{}, pkgerrors.WithMessage(err, "unable to parse a level")
	}

	parsedEncoding, err := ParseEncoding(encoding)
	if err != nil {
		return nil, zap.AtomicLevel{}, pkgerrors.WithMessage(err, "unable to parse an encoding")
	}

	config := zap.Config{
		Level:             parsedLevel,
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          parsedEncoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"app": app},
	}
	logger, err := config.Build()
	if err != nil {
		return nil, zap.AtomicLevel{}, pkgerrors.WithMessage(err, "unable to build a logger")
	}

	return logger, parsedLevel, nil
}

// NewDummy - return a new dummy logger, actually is a zap no-op Logger
func NewDummy() *zap.Logger {
	return zap.NewNop()
}
