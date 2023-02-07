package routes

import (
	"net/http"

	"github.com/jkrus/master_api/internal/io/http/routes/files_controller"
	"github.com/jkrus/master_api/internal/io/http/routes/wrappers"
)

func (r router) initFileRoutes(controller *files_controller.FileController) {
	sr := r.router.PathPrefix("/api/file").Subrouter().StrictSlash(true)

	sr.HandleFunc("/upload", wrappers.WrapJSONHandler(controller.Create)).Methods(http.MethodPost)
	sr.HandleFunc("/download", wrappers.WrapXLSHandler(controller.GetFile)).Methods(http.MethodGet)
	sr.HandleFunc("/delete", wrappers.WrapJSONHandler(controller.DeleteFile)).Methods(http.MethodDelete)

}
