package err_const

import (
	"github.com/jkrus/master_api/pkg/errors"
)

const (

	/*HTTP Statuses*/

	MsgJsonUnMarshal = "не удалось декодировать JSON"
	ErrJsonUnMarshal = errors.Const(MsgJsonUnMarshal)
	MsgJsonMarshal   = "не удалось упаковать данные в JSON"
	ErrJsonMarshal   = errors.Const(MsgJsonMarshal)
	MsgBadRequest    = "ошибка параметров запроса"
	ErrBadRequest    = errors.Const(MsgBadRequest)

	// Transport

	ErrResponseWrite = errors.Const("ошибка записи ответа")

	// DB

	ErrDatabaseRecordNotFound = errors.Const("запись не найдена")

	// MinIO

	ErrKeyDoesNotExist     = errors.Const("The specified key does not exist.")
	ErrMinioRecordNotFound = errors.Const("запись не найдена")
)
