package tracing

import (
	"github.com/streadway/amqp"

	"github.com/jkrus/master_api/pkg/errors"
)

// контейнер для извлечения тэгов трассировки для
type amqpHeadersCarrier map[string]interface{}

// ForeachKey  имплементация opentracing.TextMapReader
func (c amqpHeadersCarrier) ForeachKey(handler func(key, value string) error) error {
	for k, val := range c {
		v, ok := val.(string)
		if !ok {
			continue
		}

		if err := handler(k, v); err != nil {
			return err
		}
	}

	return nil
}

// Set имплементация opentracing.TextMapReader
func (c amqpHeadersCarrier) Set(key, value string) {
	c[key] = value
}

// InjectSpanToAMQPHeader добавляет спан в заголовки AMQP сообщения
func InjectSpanToAMQPHeader(span opentracing.Span, header amqp.Table) error {
	// если нет входящего спана - ничего не делаем
	if span == nil {
		return nil
	}

	return span.Tracer().Inject(span.Context(), opentracing.TextMap, amqpHeadersCarrier(header))
}

// StartSpanFromAMQPHeader извлекает спан из заголовков AMQP сообщения или открывает новый
func StartSpanFromAMQPHeader(operationName string, tracer opentracing.Tracer, header amqp.Table) (opentracing.Span, error) {
	spanCtx, err := tracer.Extract(opentracing.TextMap, amqpHeadersCarrier(header))
	if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
		return nil, errors.Wrap(err, "extract span from rabbit header")
	}

	span := tracer.StartSpan(operationName, opentracing.FollowsFrom(spanCtx))

	return NewSpanWrapper(span), nil
}
