package http

import (
	"fmt"
	"net/http"
	"os"
	"syscall"

	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/bl"
	"github.com/jkrus/master_api/internal/config"
	"github.com/jkrus/master_api/internal/io/http/i_http"
	"github.com/jkrus/master_api/internal/io/http/routes"
	"github.com/jkrus/master_api/pkg/shutdown"
)

const (
	MsgHTTPServerRunningError = "http server running error"
	MsgApplicationIsRunning   = "application is running"
)

// httpServer
type httpServer struct {
	logger   *zap.Logger // Логер
	bl       *bl.BL      // Бизнес логика
	finished chan bool   // Канал о завершении работы сервера
}

// NewHTTPServer конструктор c di
func NewHTTPServer(logger *zap.Logger, bl *bl.BL, finished chan bool) i_http.IHTTPServer {
	return &httpServer{
		logger:   logger,
		bl:       bl,
		finished: finished,
	}
}

// run запуск сервера
func (s httpServer) run(logger *zap.Logger, address string) error {
	settings := config.GetConfig()

	httpServer := &http.Server{
		Addr:         address,
		ReadTimeout:  settings.HTTPServerReadTimeOut,
		WriteTimeout: settings.HTTPServerWriteTimeOut,
		Handler:      routes.InitRoutes(s.logger, s.bl),
	}

	go shutdown.Graceful(logger, []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		httpServer)

	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// Run запускает го-рутину, которая стартует сервер и после его завершения посылает в канал информацию о завершении
func (s httpServer) Run(logger *zap.Logger) {
	address := getAddress()

	go func() {
		if err := s.run(logger, address); err != nil {
			s.logger.Error(MsgHTTPServerRunningError, zap.Error(err))
		}
		s.finished <- true
	}()
	s.logger.Info(MsgApplicationIsRunning, zap.String("address", address))
}

// getAddress метод получения адреса для сервера
func getAddress() string {
	settings := config.GetConfig()
	return fmt.Sprintf("%s:%s", settings.ServerHost, settings.ServerPort)
}
