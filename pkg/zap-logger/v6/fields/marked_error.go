package fields

import (
	"go.uber.org/zap"
)

// MarkedError ...
func MarkedError(key string, err error) zap.Field {
	return AddSuffix(zap.NamedError(key, err), 0, "error")
}
