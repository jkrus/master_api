package extractors

import (
	"bytes"
	"net/http"

	"go.uber.org/zap"
)

// Header ...
type Header struct {
	exceptions map[string]bool
}

// ...
var (
	DefaultHeader = NewHeader("Authorization")
)

// NewHeader ...
func NewHeader(exceptions ...string) Header {
	// reorder exceptions for the http.Header.WriteSubset() method
	exceptionSet := make(map[string]bool)
	for _, exception := range exceptions {
		exceptionSet[http.CanonicalHeaderKey(exception)] = true
	}

	return Header{exceptionSet}
}

// Extract ...
//
// It never returns an error.
//
func (extractor Header) Extract(request *http.Request) ([]zap.Field, error) {
	var buffer bytes.Buffer
	// bytes.Buffer doesn't return any errors
	request.Header.WriteSubset(&buffer, extractor.exceptions) // nolint: errcheck, gosec

	return []zap.Field{zap.String("header", buffer.String())}, nil
}
