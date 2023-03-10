package files_controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jkrus/master_api/internal/common/err_const"
	"github.com/jkrus/master_api/internal/io/http/common/response_models"
	"github.com/jkrus/master_api/internal/io/http/common/utils"
	"github.com/jkrus/master_api/internal/io/http/routes/files_controller/models"
	"github.com/jkrus/master_api/pkg/errors"
)

func (fc *FileController) Create(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("CREATE FILE")

	/*fileBody := &models.File{}
	if err := json.NewDecoder(r.Body).Decode(fileBody); err != nil {
		return nil, err
	}*/

	// TODO maximum file size
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", err
	}

	fc.logger.Debug("decode create upload fileInfo dto")
	orgID, err := utils.ExtractIdByFormRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}
	orgIDString := strconv.Itoa(int(orgID))

	orgIDString = "orgid" + orgIDString
	files, ok := r.MultipartForm.File["file"]
	if !ok || len(files) == 0 {
		return "", err_const.ErrBadRequest
	}
	userUuid, ok := r.MultipartForm.Value["UserUuid"]
	if !ok || len(userUuid) == 0 {
		return "", err_const.ErrBadRequest
	}
	statusIdStr, ok := r.MultipartForm.Value["FileStatusId"]
	if !ok || len(statusIdStr) == 0 {
		return "", err_const.ErrBadRequest
	}
	statusId, err := strconv.Atoi(statusIdStr[0])
	if err != nil {
		return "", err_const.ErrBadRequest
	}
	typeIdStr, ok := r.MultipartForm.Value["FileTypeId"]
	if !ok || len(typeIdStr) == 0 {
		return "", err_const.ErrBadRequest
	}
	typeId, err := strconv.Atoi(statusIdStr[0])
	if err != nil {
		return "", err_const.ErrBadRequest
	}

	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	file := models.File{
		UserUuid: userUuid[0],
		Name:     fileInfo.Filename,
		StatusId: uint(statusId),
		TypeId:   uint(typeId),
		Reader:   fileReader,
		Size:     fileInfo.Size,
	}

	res, err := fc.bl.FileLogic.File.Create(r.Context(), orgIDString, file.ToDTO())
	if err != nil {
		return "", err
	}
	w.WriteHeader(http.StatusCreated)

	return (&models.File{}).FromDTO(res), nil
}

func (fc *FileController) GetFile(w http.ResponseWriter, r *http.Request) ([]byte, string, error) {
	fc.logger.Info("GET FILE")

	fc.logger.Debug("get orgID from URL")
	orgID, err := utils.ExtractIdFromParamsRequest(r)
	if err != nil {
		return nil, "", errors.And(err, err_const.ErrBadRequest)
	}
	orgIDString := strconv.Itoa(int(orgID))

	orgIDString = "orgid" + orgIDString

	fc.logger.Debug("get fileUUID from context")
	fileUUID, err := utils.ExtractUUIDFromRequest(r)
	if err != nil {
		return nil, "", errors.And(err, err_const.ErrBadRequest)
	}
	res, err := fc.bl.FileLogic.File.DownloadFile(r.Context(), orgIDString, fileUUID)
	if err != nil {
		return nil, "", err
	}

	return res.Data, res.Name, nil
}

func (fc *FileController) DeleteFile(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("DELETE FILE")

	orgID, err := utils.ExtractIdFromParamsRequest(r)
	if err != nil || orgID <= 0 {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}
	orgIDString := strconv.Itoa(int(orgID))

	orgIDString = "orgid" + orgIDString

	fc.logger.Debug("get fileUUID from context")
	fileUUID, err := utils.ExtractUUIDFromParamsRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}

	err = fc.bl.FileLogic.File.Delete(r.Context(), orgIDString, fileUUID)
	if err != nil {
		return nil, err
	}

	return response_models.RestResponse{Message: "OK"}, nil
}

func (fc *FileController) UpdateFile(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	fc.logger.Info("UPDATE FILE")

	fileUuid, err := utils.ExtractUUIDFromRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}
	data := &models.File{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, err
	}

	result, err := fc.bl.FileLogic.File.Update(r.Context(), fileUuid, data.ToDTO())
	if err != nil {
		return nil, err
	}

	return (&models.File{}).FromDTO(result), nil
}
