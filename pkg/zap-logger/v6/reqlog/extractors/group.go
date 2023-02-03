package extractors

import (
	"net/http"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate mockery --name=Extractor --inpackage --case=underscore --testonly

// Extractor ...
type Extractor interface {
	Extract(request *http.Request) ([]zap.Field, error)
}

// Group ...
type Group struct {
	extractors []Extractor
}

// NewGroup ...
func NewGroup(extractors ...Extractor) Group {
	return Group{extractors}
}

// Extract ...
func (extractor Group) Extract(request *http.Request) ([]zap.Field, error) {
	var fields []zap.Field
	for _, childExtractor := range extractor.extractors {
		childFields, err := childExtractor.Extract(request)
		if err != nil {
			return nil, pkgerrors.WithMessage(err, "unable to extract fields")
		}

		fields = append(fields, childFields...)
	}

	return fields, nil
}
