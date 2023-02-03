package tracing

import (
	"context"
	"runtime"
	"strings"
	"sync/atomic"
)

var _ opentracing.Span = (*spanWrapper)(nil)

const (
	operational int32 = iota + 1
	finished
)

// nolint
var noopTracer opentracing.NoopTracer

type inheritanceFunc func(opentracing.SpanContext) opentracing.SpanReference

// обертка над opentracing.Span, позволяющая безопасно вызывать
// финиш span-а несколько раз
type spanWrapper struct {
	opentracing.Span
	State int32
}

// NewSpanWrapper конструктор обертки спана
func NewSpanWrapper(span opentracing.Span) opentracing.Span {
	return &spanWrapper{
		Span:  span,
		State: operational,
	}
}

// Finish устанавливает временную метку состояния и завершает span
func (s *spanWrapper) Finish() {
	if s == nil || s.Span == nil {
		return
	}
	if swapped := atomic.CompareAndSwapInt32(&s.State, operational, finished); swapped {
		s.Span.Finish()
	}
}

// FollowSpan возвращает следующий спан в цепочке трейсинга.
func FollowSpan(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	if operationName == "" {
		operationName = defaultOperationName()
	}
	return inheritance(ctx, opentracing.FollowsFrom, operationName, opts...)
}

// FollowCallerSpan - возвращает следующий спан в цепочке трейсинга, в качестве
// operationName подставляется имя вызывающего метода
func FollowCallerSpan(ctx context.Context, opts ...opentracing.StartSpanOption) opentracing.Span {
	operationName := defaultOperationName()
	return inheritance(ctx, opentracing.FollowsFrom, operationName, opts...)
}

// FollowWithContext возвращает дочерний спан на основе родительского полученного из контекста
func FollowWithContext(
	ctx context.Context,
	operationName string,
	opts ...opentracing.StartSpanOption,
) (context.Context, opentracing.Span) {
	if operationName == "" {
		operationName = defaultOperationName()
	}
	span := FollowSpan(ctx, operationName, opts...)
	return opentracing.ContextWithSpan(ctx, span), span
}

// FollowCallerWithContext возвращает следующий спан в цепочке трейсинга, в качестве
// operationName подставляется имя вызывающего метода, а также контекст с этим спаном
func FollowCallerWithContext(
	ctx context.Context,
	opts ...opentracing.StartSpanOption,
) (context.Context, opentracing.Span) {
	operationName := defaultOperationName()
	span := FollowSpan(ctx, operationName, opts...)
	return opentracing.ContextWithSpan(ctx, span), span
}

// Возвращает защищенный от множественного финиширования спан-потомок.
// Если родительского спана нет в контексте - создает новый из Noop-трейсера,
// либо наследует от существующего родителя.
func inheritance(
	ctx context.Context,
	inheritanceF inheritanceFunc,
	operationName string,
	opts ...opentracing.StartSpanOption,
) opentracing.Span {
	parent := opentracing.SpanFromContext(ctx)
	if parent == nil {
		return NewSpanWrapper(noopTracer.StartSpan(operationName))
	}
	opts = append(opts, inheritanceF(parent.Context()))
	child := parent.Tracer().StartSpan(operationName, opts...)
	return NewSpanWrapper(child)
}

func defaultOperationName() (fname string) {
	pc, _, _, _ := runtime.Caller(2) // nolint: gomnd, dogsled
	fname = runtime.FuncForPC(pc).Name()
	return fname[strings.Index(fname, ".")+1:]
}
