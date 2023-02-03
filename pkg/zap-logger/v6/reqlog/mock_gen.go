package reqlog

import (
	"io"
	"net/http"

	"github.com/jkrus/master_api/pkg/zap-logger/v6/reqlog/extractors"
)

//go:generate mockery --name=Extractor --inpackage --case=underscore --testonly

// Extractor is used only for generating a mock
// of the extractors.Extractor interface by the mockery tool.
type Extractor interface {
	extractors.Extractor
}

//go:generate mockery --name=ResponseWriter --inpackage --case=underscore --testonly

// ResponseWriter is used only for generating a mock
// of the http.ResponseWriter interface by the mockery tool.
type ResponseWriter interface {
	http.ResponseWriter
}

//go:generate mockery --name=Handler --inpackage --case=underscore --testonly

// Handler is used only for generating a mock
// of the http.Handler interface by the mockery tool.
type Handler interface {
	http.Handler
}

//go:generate mockery --name=Reader --inpackage --case=underscore --testonly

// Reader is used only for generating a mock
// of the io.Reader interface by the mockery tool.
type Reader interface {
	io.Reader
}
