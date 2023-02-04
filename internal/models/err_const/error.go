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

	// DB

	ErrDatabaseRecordNotFound = errors.Const("запись не найдена")
)
