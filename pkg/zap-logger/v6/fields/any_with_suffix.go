package fields

import (
	"go.uber.org/zap"
)

// AnyWithSuffix ...
func AnyWithSuffix(key string, value interface{}, suffixes ...string) zap.Field {
	// one skipped stack frame for the current function
	return AddSuffix(zap.Any(key, value), 1, suffixes...)
}
