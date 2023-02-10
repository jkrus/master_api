package files_controller

import (
	"net/http"
	"strconv"

	"github.com/jkrus/master_api/internal/common/err_const"
	"github.com/jkrus/master_api/internal/io/http/common/response_models"
	"github.com/jkrus/master_api/internal/io/http/common/utils"
	"github.com/jkrus/master_api/internal/io/http/routes/files_controller/models"
	"github.com/jkrus/master_api/pkg/errors"
)

func (f *FileController) Create(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	f.logger.Info("CREATE FILE")

	// TODO maximum file size
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", err
	}

	f.logger.Debug("decode create upload fileInfo dto")
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
	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	file := models.File{
		Name:   fileInfo.Filename,
		Size:   uint(fileInfo.Size),
		Reader: fileReader,
	}

	url, err := f.bl.File.File.Create(r.Context(), orgIDString, file.ToDTO())
	if err != nil {
		return "", err
	}
	w.WriteHeader(http.StatusCreated)

	return &models.Reference{
		Value: url,
	}, nil
}

func (f *FileController) GetFile(w http.ResponseWriter, r *http.Request) ([]byte, string, error) {
	f.logger.Info("GET FILE")

	f.logger.Debug("get orgID from URL")
	orgID, err := utils.ExtractIdFromParamsRequest(r)
	if err != nil {
		return nil, "", errors.And(err, err_const.ErrBadRequest)
	}
	orgIDString := strconv.Itoa(int(orgID))

	orgIDString = "orgid" + orgIDString

	f.logger.Debug("get fileUUID from context")
	fileUUID, err := utils.ExtractUUIDFromParamsRequest(r)
	if err != nil {
		return nil, "", errors.And(err, err_const.ErrBadRequest)
	}
	res, err := f.bl.File.File.GetFile(r.Context(), orgIDString, fileUUID)
	if err != nil {
		return nil, "", err
	}

	return res.Bytes, res.Name, nil
}

func (f *FileController) DeleteFile(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	f.logger.Info("DELETE FILE")

	orgID, err := utils.ExtractIdFromParamsRequest(r)
	if err != nil || orgID <= 0 {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}
	orgIDString := strconv.Itoa(int(orgID))

	orgIDString = "orgid" + orgIDString

	f.logger.Debug("get fileUUID from context")
	fileUUID, err := utils.ExtractUUIDFromParamsRequest(r)
	if err != nil {
		return nil, errors.And(err, err_const.ErrBadRequest)
	}

	err = f.bl.File.File.Delete(r.Context(), orgIDString, fileUUID)
	if err != nil {
		return nil, err
	}

	return response_models.RestResponse{Message: "OK"}, nil
}
