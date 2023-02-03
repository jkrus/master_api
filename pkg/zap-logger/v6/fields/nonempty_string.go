package fields

import (
	"go.uber.org/zap"
)

// NonemptyString ...
//
// It's same as zap.String but skip empty ones.
//
func NonemptyString(key string, value string) zap.Field {
	if len(value) == 0 {
		return zap.Skip()
	}

	return zap.String(key, value)
}
