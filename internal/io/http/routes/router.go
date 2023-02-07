package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/bl"
	"github.com/jkrus/master_api/internal/io/http/middlewares"
	"github.com/jkrus/master_api/internal/io/http/routes/files_controller"
	"github.com/jkrus/master_api/internal/io/http/routes/ping_controller"
)

type router struct {
	bl          *bl.BL
	router      *mux.Router
	middlewares middlewares.IMiddleWares
	logger      *zap.Logger
}

func InitRoutes(logger *zap.Logger, bl *bl.BL) http.Handler {
	r := &router{
		logger:      logger,
		bl:          bl,
		router:      mux.NewRouter(),
		middlewares: middlewares.NewMiddleWares(bl),
	}
	r.router.Use(
		r.middlewares.GetObservabilityMiddleware(logger),
	)
	r.initPingRoutes(ping_controller.NewPingController(bl))
	r.initFileRoutes(files_controller.NewFileController(logger, bl))

	return r.router
}
