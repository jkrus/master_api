package errors

import (
	"bytes"
)

var _ error = &wrappedError{}
var _ wrapped = &wrappedError{}

type wrappedError struct {
	msgs []string
	err  error

	ctx *errContext
}

// Unwrap разматывает врапнутую ошибку
func (w *wrappedError) Unwrap() error {
	return w.err
}

func (w *wrappedError) Error() string {
	var buf bytes.Buffer
	for i := len(w.msgs) - 1; i >= 0; i-- {
		buf.WriteString(w.msgs[i])
		buf.WriteString(": ")
	}
	buf.WriteString(w.err.Error())
	return buf.String()
}

// As для поиска контекста
func (w *wrappedError) As(target interface{}) bool {
	switch v := target.(type) {
	case **errContext:
		if w.ctx != nil {
			*v = w.ctx
			return true
		}
	case **wrappedError:
		*v = w
		return true
	}
	return false
}

// Reporter возвращает репорт контекста
func (w *wrappedError) Reporter() ContextReporter {
	right, ok := w.err.(Reporter)
	if !ok || right == nil {
		return w.ctx
	}

	return combinedContextReporter{
		left:  w.ctx,
		right: right.Reporter(),
	}
}

// combinedContextReporter отдача репорта из двух источников, самая являющаяся репортером
type combinedContextReporter struct {
	left  ContextReporter
	right ContextReporter
}

// Report для реализации ContextReporter
func (c combinedContextReporter) Report(builder ContextReportBuilder) {
	c.right.Report(builder)
	c.left.Report(builder)
}

// reloc релокация места возникновения ошибки на указанную позицию (файл, строка)
func (c combinedContextReporter) reloc(file string, line int) {
	if c.right != nil {
		c.right.reloc(file, line)
	} else if c.left != nil {
		c.left.reloc(file, line)
	}
}
