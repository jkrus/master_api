package files_controller

import (
	"encoding/json"
	"net/http"

	"github.com/jkrus/master_api/internal/common/err_const"
	"github.com/jkrus/master_api/internal/io/http/common/response_models"
	"github.com/jkrus/master_api/internal/io/http/common/utils"
	"github.com/jkrus/master_api/internal/io/http/routes/files_controller/models"
	"github.com/jkrus/master_api/pkg/errors"
)

func (fc *FileController) CreateFileStatus(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("CREATE FILE_STATUS")

	data := &models.FileStatus{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, err
	}

	result, err := fc.bl.FileLogic.FileStatus.Create(r.Context(), data.ToDTO())
	if err != nil {
		return nil, err
	}

	return (&models.FileStatus{}).FromDTO(result), nil
}

func (fc *FileController) GetFileStatusById(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("GET FILE_STATUS")

	fc.logger.Debug("get fileStatusId from URL")
	fileStatusId, err := utils.ExtractIdFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}

	result, err := fc.bl.FileLogic.FileStatus.GetById(r.Context(), fileStatusId)
	if err != nil {
		return nil, err
	}

	return (&models.FileStatus{}).FromDTO(result), nil
}

func (fc *FileController) UpdateFileStatus(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("UPDATE FILE_STATUS")

	fileStatusId, err := utils.ExtractIdFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}
	data := &models.FileStatus{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, err
	}

	result, err := fc.bl.FileLogic.FileStatus.Update(r.Context(), fileStatusId, data.ToDTO())
	if err != nil {
		return nil, err
	}

	return (&models.FileStatus{}).FromDTO(result), nil
}

func (fc *FileController) DeleteFileStatus(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("DELETE FILE_STATUS")

	fileStatusId, err := utils.ExtractIdFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}

	err = fc.bl.FileLogic.FileStatus.Delete(r.Context(), fileStatusId)
	if err != nil {
		return nil, err
	}

	return response_models.RestResponse{Message: "OK"}, nil
}
