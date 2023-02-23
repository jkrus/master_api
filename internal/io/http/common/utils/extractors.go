package utils

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/jkrus/master_api/pkg/errors"
)

const (
	ErrParseIdFromRequest   = "ошибка конвертации id в uint после получения из запроса"
	ErrParseUUIDFromRequest = "не верный формат uuid"
	ErrIdMissing            = "в запросе не указан идентификатор"
)

// ExtractIdFromRequest извлекает id uint из запроса
func ExtractIdFromRequest(r *http.Request) (uint, error) {
	if idString, ok := mux.Vars(r)["id"]; ok {
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			return 0, errors.Ctx().Str("id-str", idString).New(ErrParseIdFromRequest)
		}
		return uint(id), nil
	}
	return 0, errors.New(ErrIdMissing)
}

// ExtractIdByFormRequest извлекает id uint из формы. Форма должна быть распарсена зарнее
func ExtractIdByFormRequest(r *http.Request) (uint, error) {
	idString := r.Form.Get("id")
	if idString != "" {
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			return 0, errors.Ctx().Str("id-str", idString).New(ErrParseIdFromRequest)
		}
		return uint(id), nil
	}
	return 0, errors.New(ErrIdMissing)
}

// ExtractIdFromParamsRequest извлекает id uint из из параметров запроса.
func ExtractIdFromParamsRequest(r *http.Request) (uint, error) {
	idString := r.URL.Query().Get("id")
	if idString != "" {
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			return 0, errors.Ctx().Str("id-str", idString).New(ErrParseIdFromRequest)
		}
		return uint(id), nil
	}
	return 0, errors.New(ErrIdMissing)
}

// ExtractUUIDFromParamsRequest извлекает id uint из параметров запроса.
func ExtractUUIDFromParamsRequest(r *http.Request) (string, error) {
	uid := r.URL.Query().Get("uuid")
	if uid != "" {
		id, err := uuid.Parse(uid)
		if err != nil {
			return "", errors.Ctx().Str("uuid-str", uid).New(ErrParseUUIDFromRequest)
		}
		return id.String(), nil
	}
	return "", errors.New(ErrIdMissing)
}

// ExtractUUIDFromRequest - извлекает uuid string из запроса
func ExtractUUIDFromRequest(r *http.Request) (string, error) {
	if uid, ok := mux.Vars(r)["uuid"]; ok {
		id, err := uuid.Parse(uid)
		if err != nil {
			return "", errors.Ctx().Str("uuid-str", uid).New(ErrParseUUIDFromRequest)
		}
		return id.String(), nil
	}
	return "", errors.New(ErrIdMissing)
}
