package tracing

import (
	"net/http"

	"github.com/opentracing/opentracing-go/ext"

	"github.com/jkrus/master_api/pkg/errors"
)

// StartSpanFromRequest стартует спан по данным запроса, если контекст
// спана отсутсвует в запросе - создается новый родительский спан
func StartSpanFromRequest(tracer opentracing.Tracer, r *http.Request) (opentracing.Span, error) {
	span := opentracing.SpanFromContext(r.Context())

	if span == nil {
		spanCtx, err := tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header),
		)

		if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
			return nil, errors.Wrap(err, "extract span from request")
		}

		// запрашиваем новый спан у трейсера по данным контекста
		span = tracer.StartSpan(r.URL.Path, ext.RPCServerOption(spanCtx))
		ext.HTTPMethod.Set(span, r.Method)

		var rurl = r.RequestURI

		if rurl == "" {
			rurl = r.URL.String()
		}

		ext.HTTPUrl.Set(span, rurl)

		if err = InjectSpanToRequest(span, r); err != nil {
			return nil, errors.Wrap(err, "span injecting")
		}
	}

	return NewSpanWrapper(span), nil
}

// InjectSpanToRequest помещает контекст спана в http запрос
func InjectSpanToRequest(span opentracing.Span, r *http.Request) error {
	if span == nil {
		return opentracing.ErrSpanContextNotFound
	}

	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)
}
