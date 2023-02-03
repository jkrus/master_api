package errors

import (
	"errors"
)

// New обёртка для errors.New из стандартной библиотеки
func New(msg string) error {
	return Ctx().Loc(1).New(msg)
}

// Newf обёртка для fmt.Errorf из стандартной библиотеки
func Newf(format string, a ...interface{}) error {
	return Ctx().Loc(1).Newf(format, a...)
}

// Wrap аннотация ошибки данным сообщением msgs, "заматывание"
func Wrap(err error, msg string) error {
	return Ctx().Loc(1).Wrap(err, msg)
}

// Wrapf аннотация ошибки с помощью данного форматированного сообщения
func Wrapf(err error, format string, a ...interface{}) error {
	return Ctx().Loc(1).Wrapf(err, format, a...)
}

// Unwrap возвращает err.Unwrap(), если у типа имеется данный метод, иначе возвращает сам err. Т.е. эта функция
// возвращает ошибку очищенную от аннотаций
func Unwrap(err error) error {
	if res, ok := err.(wrapped); ok {
		return res.Unwrap()
	}
	return err
}

// Is сообщает, имеется ли в цепочке вложенных в err ошибок данная (target).
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As ищет первую ошибку в цепочке вложенных в err с типом и делает target, который является указателем на объект
// чей тип реализует встроенный интерфейс error. Если target не удовлетворяет данному условия, то функция паникует.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// And возвращает ошибку являющуюся списком
func And(err1, err2 error) error {
	if err2 == nil {
		return err1
	}
	if err1 == nil {
		return err2
	}
	var l1 List
	if v, ok := err1.(List); ok {
		l1 = v
	} else {
		l1 = make(List, 0, 4)
		l1 = append(l1, err1)
	}

	if v, ok := err2.(List); ok {
		l1 = append(l1, v...)
	} else {
		l1 = append(l1, err2)
	}

	return l1
}

// AsList ищет тип `errors.List` во вложенных ошибках и возвращает его. Если списка не найдено, то создаётся новый
// список, состоящий из текущей ошибки.
func AsList(err error) List {
	// Если текущая ошибка является декорацией к списку (wrappedError(lst)), то возвращается список ошибок, где к
	// каждой применена текущая декорация. Это костыль, который сломается при декорировании любым другим способом,
	// например, с помощью `fmt.Errorf("wrapper message: %w", err)`, но в нашей практике его будет достаточно
	if v, ok := err.(*wrappedError); ok {
		if lst, ok := v.err.(List); ok {
			newLst := make(List, len(lst))
			for i, e := range lst {
				newLst[i] = &wrappedError{
					msgs: v.msgs,
					err:  e,
					ctx:  v.ctx,
				}
			}
			return newLst
		}
	}

	var lst List
	if !As(err, &lst) {
		return List{err}
	}

	return lst
}

// Reloc изменяет место возникновения ошибки если оно установлено или добавляет его.
// имеет смысл, когда фактическое место возникновения не информативно.
func Reloc(err error) error {
	report := GetContextReporter(err)
	if report == nil {
		return Ctx().Loc(1).Just(err)
	}

	ctx := Ctx().Loc(1).(*errContext)
	report.reloc(ctx.loc.file, ctx.loc.line)

	return err
}

// GetContextReporter получает репортера из ошибки
func GetContextReporter(err error) ContextReporter {
	v, ok := err.(Reporter)
	if !ok || v == nil {
		return nil
	}

	return v.Reporter()
}
