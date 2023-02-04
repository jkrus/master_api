package routes

import (
	"net/http"

	"github.com/jkrus/master_api/internal/io/http/routes/ping_controller"
	"github.com/jkrus/master_api/internal/io/http/routes/wrappers"
)

func (r router) initPingRoutes(controller *ping_controller.PingController) {
	sr := r.router.PathPrefix("/api").Subrouter().StrictSlash(true)

	sr.HandleFunc("/ping", wrappers.WrapJSONHandler(controller.Ping)).Methods(http.MethodGet)

}
