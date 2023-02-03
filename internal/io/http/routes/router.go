package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jkrus/master_api/internal/bl"
)

type router struct {
	bl     *bl.BL
	router *mux.Router
}

func InitRoutes(bl *bl.BL) http.Handler {
	r := &router{
		bl:     bl,
		router: mux.NewRouter(),
	}
	return r.router
}
