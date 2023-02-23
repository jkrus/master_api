package routes

import (
	"net/http"

	"github.com/jkrus/master_api/internal/io/http/routes/users_controller"
	"github.com/jkrus/master_api/internal/io/http/routes/wrappers"
)

func (r router) initUsersRoutes(controller *users_controller.UserController) {
	sr := r.router.PathPrefix("/api").Subrouter().StrictSlash(true)

	/* Users */

	sr.HandleFunc("/users", wrappers.WrapJSONHandler(controller.CreateUser)).Methods(http.MethodPost)
	sr.HandleFunc("/users/{uuid}", wrappers.WrapJSONHandler(controller.GetUserByUuid)).Methods(http.MethodGet)
	sr.HandleFunc("/users/{uuid}", wrappers.WrapJSONHandler(controller.UpdateUser)).Methods(http.MethodPut)
	sr.HandleFunc("/users/{uuid}", wrappers.WrapJSONHandler(controller.DeleteUser)).Methods(http.MethodDelete)

}
