package zaplogger

import (
	"errors"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"
)

// ParseLevel ...
func ParseLevel(level string) (zap.AtomicLevel, error) {
	if level == "" {
		return zap.AtomicLevel{}, errors.New("empty level")
	}

	atomicLevel := zap.NewAtomicLevel()
	if err := atomicLevel.UnmarshalText([]byte(strings.ToLower(level))); err != nil {
		return zap.AtomicLevel{}, pkgerrors.WithMessagef(err, "unknown level %s", level)
	}

	return atomicLevel, nil
}
