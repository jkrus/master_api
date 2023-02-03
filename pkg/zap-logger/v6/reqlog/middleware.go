package reqlog

import (
	"crypto/rand"
	"io"
	"net/http"

	"go.uber.org/zap"

	zaplogger "github.com/jkrus/master_api/pkg/zap-logger/v6"
	"github.com/jkrus/master_api/pkg/zap-logger/v6/fields"
	"github.com/jkrus/master_api/pkg/zap-logger/v6/reqlog/extractors"
)

//go:generate mockery --name=ErrHandler --inpackage --case=underscore --testonly

// ErrHandler ...
type ErrHandler interface {
	ServeHTTPErr(writer http.ResponseWriter, request *http.Request, err error)
}

// Middleware ...
type Middleware struct {
	logger       *zap.Logger
	extractor    extractors.Extractor
	errHandler   ErrHandler
	randomReader io.Reader
}

// NewMiddleware ...
func NewMiddleware(
	logger *zap.Logger,
	extractor extractors.Extractor,
	errHandler ErrHandler,
	randomReader io.Reader,
) Middleware {
	return Middleware{logger, extractor, errHandler, randomReader}
}

// NewDefaultMiddleware ...
func NewDefaultMiddleware(logger *zap.Logger) Middleware {
	// these extractors never return an error, so we can don't set an error handler
	extractor := extractors.NewGroup(extractors.DefaultRequest, extractors.DefaultHeader)
	return Middleware{logger, extractor, nil, rand.Reader}
}

func (middleware Middleware) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
	nextHandler http.HandlerFunc,
) {
	defaultFields, err := middleware.extractor.Extract(request)
	if err != nil {
		middleware.errHandler.ServeHTTPErr(writer, request, err)
		return
	}

	hash := fields.Hash("hash", fields.WithRandomSource(middleware.randomReader))
	defaultFields = append(defaultFields, hash)

	logger := middleware.logger.With(defaultFields...)
	request = AddToRequest(request, logger)
	defer zaplogger.Recover(logger)

	nextHandler(writer, request)
}
