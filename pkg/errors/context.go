package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Ctx возврат контекста ошибки
func Ctx() ContextBuilder {
	return &errContext{}
}

var _ error = &errContext{}
var _ ContextBuilder = &errContext{}

// позиция ошибки
type location struct {
	file string
	line int
}

// сущность помещаемая в контекст ошибки
type contextItem struct {
	name  string
	value itemValue
}

// контекст ошибки
type errContext struct {
	items []contextItem
	loc   *location
}

func (c *errContext) Error() string { return "" }

func (c *errContext) reloc(file string, line int) {
	if c.loc == nil {
		c.loc = &location{}
	}
	c.loc.file = file
	c.loc.line = line
}

func (c *errContext) Just(err error) error {
	switch v := err.(type) {
	case *wrappedError:
		v.ctx.items = append(v.ctx.items, c.items...)
		return v
	default:
		if GetContextReporter(err) != nil {
			c.loc = nil
		} else if c.loc == nil {
			c.Loc(1)
		}
		return &wrappedError{
			err: err,
			ctx: c,
		}
	}
}

func (c *errContext) New(msg string) error {
	if c.loc == nil {
		c.Loc(1)
	}
	return &wrappedError{
		err: errors.New(msg),
		ctx: c,
	}
}

func (c *errContext) Newf(format string, a ...interface{}) error {
	if c.loc == nil {
		c.Loc(1)
	}
	return &wrappedError{
		err: fmt.Errorf(format, a...),
		ctx: c,
	}
}

func (c *errContext) Wrap(err error, msg string) error {
	if err == nil {
		panic("wrapping nil error")
	}
	switch we := err.(type) {
	case *wrappedError:
		we.msgs = append(we.msgs, msg)
		we.ctx.items = append(we.ctx.items, c.items...)
		return we
	default:
		msgs := make([]string, 1, 4)
		msgs[0] = msg
		if GetContextReporter(err) != nil {
			// позиция уже есть, вторая не нужна
			c.loc = nil
		} else if c.loc == nil {
			c.Loc(1)
		}
		return &wrappedError{
			msgs: msgs,
			err:  err,
			ctx:  c,
		}
	}
}

func (c *errContext) Wrapf(err error, format string, a ...interface{}) error {
	if GetContextReporter(err) != nil {
		c.loc = nil
	} else if c.loc == nil {
		if c.loc == nil {
			c.Loc(1)
		}
	}
	return c.Wrap(err, fmt.Sprintf(format, a...))
}

func (c *errContext) Bool(name string, value bool) ContextBuilder {
	c.items = append(
		c.items,
		contextItem{
			name:  name,
			value: boolValue(value),
		},
	)
	return c
}

func (c *errContext) Int(name string, value int) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: intValue(value),
	})
	return c
}

func (c *errContext) Uint(name string, value uint) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uintValue(value),
	})
	return c
}

func (c *errContext) Int8(name string, value int8) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: int8Value(value),
	})
	return c
}

func (c *errContext) Int16(name string, value int16) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: int16Value(value),
	})
	return c
}

func (c *errContext) Int32(name string, value int32) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: int32Value(value),
	})
	return c
}

func (c *errContext) Int64(name string, value int64) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: int64Value(value),
	})
	return c
}

func (c *errContext) Uint8(name string, value uint8) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uint8Value(value),
	})
	return c
}

func (c *errContext) Uint16(name string, value uint16) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uint16Value(value),
	})
	return c
}

func (c *errContext) Uint32(name string, value uint32) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uint32Value(value),
	})
	return c
}

func (c *errContext) Uint64(name string, value uint64) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: uint64Value(value),
	})
	return c
}

func (c *errContext) Float32(name string, value float32) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: float32Value(value),
	})
	return c
}

func (c *errContext) Float64(name string, value float64) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: float64Value(value),
	})
	return c
}

func (c *errContext) String(name string, value string) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: stringValue(value),
	})
	return c
}

// Str синоним для String
func (c *errContext) Str(name string, value string) ContextBuilder {
	return c.String(name, value)
}

func (c *errContext) Stringer(name string, value fmt.Stringer) ContextBuilder {
	return c.String(name, value.String())
}

func (c *errContext) Strings(name string, values []string) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: stringsValue(values),
	})
	return c
}

func (c *errContext) Any(name string, value interface{}) ContextBuilder {
	c.items = append(c.items, contextItem{
		name:  name,
		value: anyValue{val: value},
	})
	return c
}

func (c *errContext) Line() ContextBuilder {
	return c.Loc(0)
}

func (c *errContext) Loc(depth int) ContextBuilder {
	_, fn, line, _ := runtime.Caller(1 + depth)
	if c.loc == nil {
		c.loc = &location{}
	}
	// оставляем только имя пакета
	if pos := strings.Index(fn, "/src/"); pos >= 0 {
		fn = fn[pos+5:]
	} else if pos = strings.Index(fn, "/pkg/mod/"); pos >= 0 {
		fn = fn[pos+9:]
	}
	c.loc.file = fn
	c.loc.line = line
	return c
}

func (c *errContext) Report(dest ContextReportBuilder) {
	if c.loc != nil {
		dest.String("err-location", fmt.Sprintf("%s:%d", c.loc.file, c.loc.line))
	}
	for _, item := range c.items {
		item.value.reportItem(item.name, dest)
	}
}
