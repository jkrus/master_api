package wrappers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jkrus/master_api/internal/models/err_const"

	"github.com/jkrus/master_api/internal/io/http/routes/common_client_models"

	"go.uber.org/zap"

	"github.com/jkrus/master_api/pkg/errors"
	zaplogger "github.com/jkrus/master_api/pkg/zap-logger/v6"
	"github.com/jkrus/master_api/pkg/zap-logger/v6/reqlog"
)

const (
	MsgRequestError       = "ошибка выполнения запроса"
	MsgJsonMarshalError   = "ошибка конвертации данных в JSON"
	MsgResponseWriteError = "ошибка записи ответа клиенту"
)

type HandlerFunc = func(w http.ResponseWriter, req *http.Request) (interface{}, error)

func WrapJSONHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := reqlog.GetFromRequest(r)
		defer func() {
			if err := logger.Sync(); err != nil {
				return
			}
		}()

		res, err := handler(w, r)
		if err != nil {
			writeErrorResponseWeb(w, r, err)
			return
		}
		resBytes, err := json.Marshal(res)
		if err != nil {
			logger.Error(MsgJsonMarshalError, zap.Error(err))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		func() {
			if _, err := w.Write(resBytes); err != nil {
				logger.Error(MsgResponseWriteError, zap.Error(err))
			}
		}()
	}
}
func writeErrorResponseWeb(w http.ResponseWriter, req *http.Request, err error) {
	reqlog.GetFromRequest(req).With(zaplogger.ExtractErrCtx(err)...).Error(MsgRequestError)

	restResponse := &common_client_models.RestResponse{
		Details: err.Error(),
	}
	var statusCode int
	statusCode, restResponse.Message = statusCodeAndErrorMessage(err)

	resBytes, _ := json.Marshal(restResponse)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(resBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func statusCodeAndErrorMessage(err error) (int, string) {
	switch {
	case errors.Is(err, err_const.ErrJsonUnMarshal):
		return http.StatusBadRequest, err_const.ErrJsonUnMarshal.Error()
	case errors.Is(err, err_const.ErrJsonMarshal):
		return http.StatusInternalServerError, err_const.ErrJsonMarshal.Error()
	case errors.Is(err, err_const.ErrDatabaseRecordNotFound):
		return http.StatusNotFound, err_const.ErrDatabaseRecordNotFound.Error()
	case strings.HasSuffix(err.Error(), "record not found"):
		return http.StatusNotFound, "Запись не найдена"
	default:
		return http.StatusInternalServerError, errors.Ctx().Wrap(err, MsgRequestError).Error()
	}
}
