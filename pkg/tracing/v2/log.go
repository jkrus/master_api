package tracing

import (
	"github.com/opentracing/opentracing-go/log"

	"github.com/jkrus/master_api/pkg/errors"
)

type errCtxExtractor struct {
	fields []log.Field
}

// LogError - устанавливает событие ошибки в спане
func LogError(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogFields(extractErrCtx(err)...)
}

func extractErrCtx(err error) (fields []log.Field) {
	if errCtx := errors.GetContextReporter(err); errCtx != nil {
		builder := &errCtxExtractor{}
		errCtx.Report(builder)

		fields = builder.Fields()
	}

	fields = append(fields, log.Error(err))

	return fields
}

// Fields - возвращает извлеченные поля логгера из контекста ошибки
func (fe *errCtxExtractor) Fields() []log.Field {
	return fe.fields
}

func (fe *errCtxExtractor) Bool(name string, value bool) {
	fe.fields = append(fe.fields, log.Bool(name, value))
}

func (fe *errCtxExtractor) Int(name string, value int) {
	fe.fields = append(fe.fields, log.Int(name, value))
}

func (fe *errCtxExtractor) Int8(name string, value int8) {
	fe.fields = append(fe.fields, log.Int32(name, int32(value)))
}

func (fe *errCtxExtractor) Int16(name string, value int16) {
	fe.fields = append(fe.fields, log.Int32(name, int32(value)))
}

func (fe *errCtxExtractor) Int32(name string, value int32) {
	fe.fields = append(fe.fields, log.Int32(name, value))
}

func (fe *errCtxExtractor) Int64(name string, value int64) {
	fe.fields = append(fe.fields, log.Int64(name, value))
}

func (fe *errCtxExtractor) Uint(name string, value uint) {
	fe.fields = append(fe.fields, log.Uint64(name, uint64(value)))
}

func (fe *errCtxExtractor) Uint8(name string, value uint8) {
	fe.fields = append(fe.fields, log.Uint32(name, uint32(value)))
}

func (fe *errCtxExtractor) Uint16(name string, value uint16) {
	fe.fields = append(fe.fields, log.Uint32(name, uint32(value)))
}

func (fe *errCtxExtractor) Uint32(name string, value uint32) {
	fe.fields = append(fe.fields, log.Uint32(name, value))
}

func (fe *errCtxExtractor) Uint64(name string, value uint64) {
	fe.fields = append(fe.fields, log.Uint64(name, value))
}

func (fe *errCtxExtractor) Float32(name string, value float32) {
	fe.fields = append(fe.fields, log.Float32(name, value))
}

func (fe *errCtxExtractor) Float64(name string, value float64) {
	fe.fields = append(fe.fields, log.Float64(name, value))
}

func (fe *errCtxExtractor) String(name, value string) {
	fe.fields = append(fe.fields, log.String(name, value))
}

func (fe *errCtxExtractor) Strings(name string, values []string) {
	fe.fields = append(fe.fields, log.Object(name, values))
}

func (fe *errCtxExtractor) Any(name string, value interface{}) {
	fe.fields = append(fe.fields, log.Object(name, value))
}

func (fe *errCtxExtractor) Err(err error) {
	fe.fields = append(fe.fields, log.Error(err))
}
