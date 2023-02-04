package reqlog

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/jkrus/master_api/pkg/zap-logger/v6/ctxlog"
)

// GetFromRequest ...
func GetFromRequest(request *http.Request) *zap.Logger {
	return ctxlog.GetFromCtx(request.Context())
}

// AddToRequest ...
func AddToRequest(request *http.Request, logger *zap.Logger) *http.Request {
	return request.WithContext(ctxlog.AddToCtx(request.Context(), logger))
}

// WithFields ...
func WithFields(request *http.Request, fields ...zap.Field) *http.Request {
	return AddToRequest(request, GetFromRequest(request).With(fields...))
}
