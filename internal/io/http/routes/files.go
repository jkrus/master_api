package routes

import (
	"net/http"

	"github.com/jkrus/master_api/internal/io/http/routes/files_controller"
	"github.com/jkrus/master_api/internal/io/http/routes/wrappers"
)

func (r router) initFileRoutes(controller *files_controller.FileController) {
	sr := r.router.PathPrefix("/api").Subrouter().StrictSlash(true)

	sr.HandleFunc("/file/upload", wrappers.WrapJSONHandler(controller.Create)).Methods(http.MethodPost)
	sr.HandleFunc("/file/download", wrappers.WrapXLSHandler(controller.GetFile)).Methods(http.MethodGet)
	sr.HandleFunc("/file/delete", wrappers.WrapJSONHandler(controller.DeleteFile)).Methods(http.MethodDelete)

	/* File Statuses */

	sr.HandleFunc("/file_statuses", wrappers.WrapJSONHandler(controller.CreateFileStatus)).Methods(http.MethodPost)
	sr.HandleFunc("/file_statuses/{id}", wrappers.WrapJSONHandler(controller.GetFileStatusById)).Methods(http.MethodGet)
	sr.HandleFunc("/file_statuses/{id}", wrappers.WrapJSONHandler(controller.UpdateFileStatus)).Methods(http.MethodPut)
	sr.HandleFunc("/file_statuses/{id}", wrappers.WrapJSONHandler(controller.DeleteFileStatus)).Methods(http.MethodDelete)

}
