package extractors

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/jkrus/master_api/pkg/zap-logger/v6/fields"
)

// Request ...
type Request struct{}

// ...
var (
	DefaultRequest = Request{}
)

// Extract ...
//
// It never returns an error.
//
func (extractor Request) Extract(request *http.Request) ([]zap.Field, error) {
	fields := []zap.Field{
		zap.String("method", request.Method),
		zap.String("route", request.URL.Path),
		fields.Parameters("query", request.URL.Query()),
		fields.NonemptyString("rawQuery", request.URL.RawQuery),
		zap.String("host", request.Host),
		zap.String("remote", request.RemoteAddr),
	}
	return fields, nil
}
