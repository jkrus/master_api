package zaplogger

import (
	"go.uber.org/zap"

	"github.com/jkrus/master_api/pkg/errors"
)

type errCtxExtractor struct {
	fields []zap.Field
}

// NewErrCtxExtractor возвращает реализацию errors.ContextReportBuilder для
// извлечения полей логгера из контекста ошибки
func NewErrCtxExtractor() errors.ContextReportBuilder {
	return &errCtxExtractor{}
}

// ExtractErrCtx извлекает данные из контекста ошибки и возвращает срез полей
// логгера
func ExtractErrCtx(err error) (fields []zap.Field) {
	if errCtx := errors.GetContextReporter(err); errCtx != nil {
		builder := &errCtxExtractor{}
		errCtx.Report(builder)

		fields = builder.Fields()
	}

	fields = append(fields, zap.Error(err))
	return fields
}

// Fields возвращает извлечённые поля логгера из контекста ошибки
func (fe *errCtxExtractor) Fields() []zap.Field {
	return fe.fields
}

func (fe *errCtxExtractor) Bool(name string, value bool) {
	fe.fields = append(fe.fields, zap.Bool(name, value))
}

func (fe *errCtxExtractor) Int(name string, value int) {
	fe.fields = append(fe.fields, zap.Int(name, value))
}

func (fe *errCtxExtractor) Int8(name string, value int8) {
	fe.fields = append(fe.fields, zap.Int8(name, value))
}

func (fe *errCtxExtractor) Int16(name string, value int16) {
	fe.fields = append(fe.fields, zap.Int16(name, value))
}

func (fe *errCtxExtractor) Int32(name string, value int32) {
	fe.fields = append(fe.fields, zap.Int32(name, value))
}

func (fe *errCtxExtractor) Int64(name string, value int64) {
	fe.fields = append(fe.fields, zap.Int64(name, value))
}

func (fe *errCtxExtractor) Uint(name string, value uint) {
	fe.fields = append(fe.fields, zap.Uint(name, value))
}

func (fe *errCtxExtractor) Uint8(name string, value uint8) {
	fe.fields = append(fe.fields, zap.Uint8(name, value))
}

func (fe *errCtxExtractor) Uint16(name string, value uint16) {
	fe.fields = append(fe.fields, zap.Uint16(name, value))
}

func (fe *errCtxExtractor) Uint32(name string, value uint32) {
	fe.fields = append(fe.fields, zap.Uint32(name, value))
}

func (fe *errCtxExtractor) Uint64(name string, value uint64) {
	fe.fields = append(fe.fields, zap.Uint64(name, value))
}

func (fe *errCtxExtractor) Float32(name string, value float32) {
	fe.fields = append(fe.fields, zap.Float32(name, value))
}

func (fe *errCtxExtractor) Float64(name string, value float64) {
	fe.fields = append(fe.fields, zap.Float64(name, value))
}

func (fe *errCtxExtractor) String(name string, value string) {
	fe.fields = append(fe.fields, zap.String(name, value))
}

func (fe *errCtxExtractor) Strings(name string, values []string) {
	fe.fields = append(fe.fields, zap.Strings(name, values))
}

func (fe *errCtxExtractor) Any(name string, value interface{}) {
	fe.fields = append(fe.fields, zap.Any(name, value))
}

func (fe *errCtxExtractor) Err(err error) {
	fe.fields = append(fe.fields, zap.Error(err))
}
