# errors

функционал для работы с ошибками

## Оборачивание (аннотирование) ошибок

осуществляется с помощью функций `Wrap` и `Wrapf`:

```go
return errors.Wrapf(err, "parse config file %s", fileName)
```

## Определение ошибок

Осуществляется с помощью функций `Is` и `As` – обертки функций из
стандартной библиотеки

## Контексты ошибок.
В ошибки может быть вшит структурированный контекст, который затем можно,
передавать через параметр контекста и раскрывать, например в логгере. 

Пример:
```go
err := errors.Ctx().String("id", "845aa443-7a63-439d-950b-4b93d72da903").String("ipaddr", "127.0.0.1").New("some error")
``` 

поддерживаются следующие методы сохранения значений:
* `Int(name string, value int)`
* `Int8(name string, value int8)`
* `Int16(name string, value int16)`
* `Int32(name string, value int32)`
* `Int64(name string, value int64)`
* `Uint(name string, value uint)`
* `Uint8(name string, value uint8)`
* `Uint16(name string, value uint16)`
* `Uint32(name string, value uint32)`
* `Uint64(name string, value uint64)`
* `Float32(name string, value float32)`
* `Float64(name string, value float64)`
* `String(name string, value string)`
* `Strings(name string, values []string)`
* `Any(name string, value interface{})`
* `Line()` 

Метод `Line()` в данном применении, сохраняет позицию `<файл>:<строка>` на которой
произошёл его вызов. Его вызов заменяет предыдущее значение (предыдущий вызов `Line`)

После сохранения необходимых контекстных значений нужно вызвать один из методов
контекста:
* `New`
* `Newf`
* `Wrap`
* `Wrapf`

для получения ошибки аннотированной контекстом

TODO: добавить работу со списками ошибок и error-group для упрощения поднятия наверх
всех нужных ошибок, например в асинхронщине.