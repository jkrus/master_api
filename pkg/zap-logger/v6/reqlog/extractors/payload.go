package extractors

import (
	"bytes"
	"io/ioutil"
	"net/http"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/pkg/zap-logger/v6/fields"
)

// Payload ...
type Payload struct {
	options []fields.PayloadOption
}

// ...
var (
	DefaultPayload = NewPayload()
)

// NewPayload ...
func NewPayload(options ...fields.PayloadOption) Payload {
	return Payload{options}
}

// Extract ...
func (extractor Payload) Extract(request *http.Request) ([]zap.Field, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "unable to read a payload")
	}
	// no a body or empty one
	if len(body) == 0 {
		return nil, nil
	}

	// restore the body, so others can read it
	request.Body = ioutil.NopCloser(bytes.NewReader(body))

	options := append([]fields.PayloadOption(nil), extractor.options...)
	options = append(options, fields.WithRequestURLSuffix(request))

	return fields.Payload("payload", body, options...), nil
}
