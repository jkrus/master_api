package routes

import (
	"net/http"

	"github.com/jkrus/master_api/internal/io/http/routes/files_controller"
	"github.com/jkrus/master_api/internal/io/http/routes/wrappers"
)

func (r router) initFileRoutes(controller *files_controller.FileController) {
	sr := r.router.PathPrefix("/api").Subrouter().StrictSlash(true)

	sr.HandleFunc("/file/upload", wrappers.WrapJSONHandler(controller.Create)).Methods(http.MethodPost)
	sr.HandleFunc("/file/download/{uuid}", wrappers.WrapXLSHandler(controller.GetFile)).Methods(http.MethodGet)
	sr.HandleFunc("/file/delete/{uuid}", wrappers.WrapJSONHandler(controller.DeleteFile)).Methods(http.MethodDelete)
	sr.HandleFunc("/file/{uuid}", wrappers.WrapJSONHandler(controller.UpdateFile)).Methods(http.MethodPut)

	/* FileStatuses */

	sr.HandleFunc("/file_statuses", wrappers.WrapJSONHandler(controller.CreateFileStatus)).Methods(http.MethodPost)
	sr.HandleFunc("/file_statuses/{id}", wrappers.WrapJSONHandler(controller.GetFileStatusById)).Methods(http.MethodGet)
	sr.HandleFunc("/file_statuses/{id}", wrappers.WrapJSONHandler(controller.UpdateFileStatus)).Methods(http.MethodPut)
	sr.HandleFunc("/file_statuses/{id}", wrappers.WrapJSONHandler(controller.DeleteFileStatus)).Methods(http.MethodDelete)

	/* FileTypes */

	sr.HandleFunc("/file_types", wrappers.WrapJSONHandler(controller.CreateFileType)).Methods(http.MethodPost)
	sr.HandleFunc("/file_types/{id}", wrappers.WrapJSONHandler(controller.GetFileTypeById)).Methods(http.MethodGet)
	sr.HandleFunc("/file_types/{id}", wrappers.WrapJSONHandler(controller.UpdateFileType)).Methods(http.MethodPut)
	sr.HandleFunc("/file_types/{id}", wrappers.WrapJSONHandler(controller.DeleteFileType)).Methods(http.MethodDelete)

}
