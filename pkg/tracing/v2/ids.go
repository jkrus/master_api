package tracing

import (
	"github.com/google/uuid"
)

// TraceIDFromSpan извлекает TraceID span'а
func TraceIDFromSpan(span opentracing.Span) string {
	return extractID(span, func(ctx jaeger.SpanContext) string {
		return ctx.TraceID().String()
	})
}

// SpanIDFromSpan извлекает SpanID span'а
func SpanIDFromSpan(span opentracing.Span) string {
	return extractID(span, func(ctx jaeger.SpanContext) string {
		return ctx.SpanID().String()
	})
}

type extractor func(jaeger.SpanContext) string

func extractID(span opentracing.Span, extract extractor) string {
	spanCtx, ok := span.Context().(jaeger.SpanContext)
	if ok {
		return extract(spanCtx)
	}

	// tracer is disabled -> span=noopSpan, span.Context=noopSpanContext -> generate uuid
	return uuid.NewString()
}
