package errors

import (
	"fmt"
)

// ContextBuilder заполнение контекста информацией и функции завершения построения
type ContextBuilder interface {
	Just(err error) error
	New(msg string) error
	Newf(format string, a ...interface{}) error
	Wrap(err error, msg string) error
	Wrapf(err error, format string, a ...interface{}) error

	Bool(name string, value bool) ContextBuilder
	Int(name string, value int) ContextBuilder
	Int8(name string, value int8) ContextBuilder
	Int16(name string, value int16) ContextBuilder
	Int32(name string, value int32) ContextBuilder
	Int64(name string, value int64) ContextBuilder
	Uint(name string, value uint) ContextBuilder
	Uint8(name string, value uint8) ContextBuilder
	Uint16(name string, value uint16) ContextBuilder
	Uint32(name string, value uint32) ContextBuilder
	Uint64(name string, value uint64) ContextBuilder
	Float32(name string, value float32) ContextBuilder
	Float64(name string, value float64) ContextBuilder
	// Deprecated: нужно использовать Str
	String(name string, value string) ContextBuilder
	Str(name string, value string) ContextBuilder
	Stringer(name string, value fmt.Stringer) ContextBuilder
	Strings(name string, values []string) ContextBuilder
	Any(name string, value interface{}) ContextBuilder
	Line() ContextBuilder
	// Loc аналогично Line, но можно указать элемент стека, откуда нужно брать позицию
	Loc(depth int) ContextBuilder

	Report(dest ContextReportBuilder)
}

// ContextReportBuilder абстракция занесения сущностей сохранённых ContextBuilder в контекст, например, логгера
type ContextReportBuilder interface {
	Bool(name string, value bool)
	Int(name string, value int)
	Int8(name string, value int8)
	Int16(name string, value int16)
	Int32(name string, value int32)
	Int64(name string, value int64)
	Uint(name string, value uint)
	Uint8(name string, value uint8)
	Uint16(name string, value uint16)
	Uint32(name string, value uint32)
	Uint64(name string, value uint64)
	Float32(name string, value float32)
	Float64(name string, value float64)
	String(name string, value string)
	Strings(name string, values []string)
	Any(name string, value interface{})
	Err(err error)
}

// Reporter описание сущностей возвращающих репортер контекста
type Reporter interface {
	Reporter() ContextReporter
}

// ContextReporter сущность для абстракции отдачи значений контекста внешним потребителя реализующим
// ContextReportBuilder
type ContextReporter interface {
	Report(builder ContextReportBuilder)
	reloc(file string, line int)
}

type wrapped interface {
	Unwrap() error
}

type itemValue interface {
	reportItem(name string, builder ContextReportBuilder)
}
