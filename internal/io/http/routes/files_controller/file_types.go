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

func (fc *FileController) CreateFileType(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("CREATE FILE_TYPE")

	data := &models.FileType{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, err
	}

	result, err := fc.bl.FileLogic.FileType.Create(r.Context(), data.ToDTO())
	if err != nil {
		return nil, err
	}

	return (&models.FileType{}).FromDTO(result), nil
}

func (fc *FileController) GetFileTypeById(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("GET FILE_TYPE")

	fc.logger.Debug("get fileTypeId from URL")
	fileTypeId, err := utils.ExtractIdFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}

	result, err := fc.bl.FileLogic.FileType.GetById(r.Context(), fileTypeId)
	if err != nil {
		return nil, err
	}

	return (&models.FileType{}).FromDTO(result), nil
}

func (fc *FileController) UpdateFileType(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("UPDATE FILE_TYPE")

	fileTypeId, err := utils.ExtractIdFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}
	data := &models.FileType{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, err
	}

	result, err := fc.bl.FileLogic.FileType.Update(r.Context(), fileTypeId, data.ToDTO())
	if err != nil {
		return nil, err
	}

	return (&models.FileType{}).FromDTO(result), nil
}

func (fc *FileController) DeleteFileType(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("DELETE FILE_TYPE")

	fileTypeId, err := utils.ExtractIdFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}

	err = fc.bl.FileLogic.FileType.Delete(r.Context(), fileTypeId)
	if err != nil {
		return nil, err
	}

	return response_models.RestResponse{Message: "OK"}, nil
}
