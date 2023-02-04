package errors

import (
	"bytes"
)

var _ error = List{}

// List тип представляющий список ошибок
type List []error

func (l List) Error() string {
	switch len(l) {
	case 0:
		// специально вызываем панику – так положено, ибо нужно явно проверять
		var err error
		return err.Error()
	case 1:
		return l[0].Error()
	default:
		var buf bytes.Buffer
		for i, err := range l {
			if i > 0 {
				buf.WriteString("; ")
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// As применяет функцию errors.As последовательно к элементам списка с ошибками, останавливаясь на успехе
func (l List) As(target interface{}) bool {
	for _, err := range l {
		if As(err, target) {
			return true
		}
	}
	return false
}

// Is применяет функцию errors.Is последовательно к элементам списка с ошибками, останавливаясь на успехе
func (l List) Is(err error) bool {
	for _, e := range l {
		if Is(e, err) {
			return true
		}
	}
	return false
}
