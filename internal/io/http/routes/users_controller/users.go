package users_controller

import (
	"encoding/json"
	"net/http"

	"github.com/jkrus/master_api/internal/common/err_const"
	"github.com/jkrus/master_api/internal/io/http/common/response_models"
	"github.com/jkrus/master_api/internal/io/http/common/utils"
	"github.com/jkrus/master_api/internal/io/http/routes/users_controller/models"
	"github.com/jkrus/master_api/pkg/errors"
)

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uc.logger.Info("CREATE USER")

	data := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, err
	}

	result, err := uc.bl.UsersLogic.Users.Create(r.Context(), data.ToDTO())
	if err != nil {
		return nil, err
	}

	return (&models.User{}).FromDTO(result), nil
}

func (uc *UserController) GetUserByUuid(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uc.logger.Info("GET USER")

	uc.logger.Debug("get userId from URL")
	userId, err := utils.ExtractUUIDFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}

	result, err := uc.bl.UsersLogic.Users.GetByUuid(r.Context(), userId)
	if err != nil {
		return nil, err
	}

	return (&models.User{}).FromDTO(result), nil
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uc.logger.Info("UPDATE USER")

	userId, err := utils.ExtractUUIDFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}
	data := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, err
	}

	result, err := uc.bl.UsersLogic.Users.Update(r.Context(), userId, data.ToDTO())
	if err != nil {
		return nil, err
	}

	return (&models.User{}).FromDTO(result), nil
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	uc.logger.Info("DELETE USER")

	userId, err := utils.ExtractUUIDFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}

	err = uc.bl.UsersLogic.Users.Delete(r.Context(), userId)
	if err != nil {
		return nil, err
	}

	return response_models.RestResponse{Message: "OK"}, nil
}
