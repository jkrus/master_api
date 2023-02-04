package main

import (
	"github.com/jkrus/master_api/internal/config"
	"github.com/jkrus/master_api/pkg/errors"
	zaplogger "github.com/jkrus/master_api/pkg/zap-logger/v6"

	"github.com/jkrus/master_api/internal"
	"github.com/jkrus/master_api/internal/bl"
	http2 "github.com/jkrus/master_api/internal/io/http"
)

func main() {
	appName := "server"

	// create logger
	logger, loggerLevel, err := zaplogger.New(appName, "info", "json")
	defer zaplogger.Recover(logger)
	if err != nil {
		logger.Error("logger creating error", zaplogger.ExtractErrCtx(errors.Ctx().Loc(2).Just(err))...)
		return
	}

	settings := config.GetConfig()
	err = loggerLevel.UnmarshalText([]byte(settings.LogLevel))
	if err != nil {
		logger.Error("can't change log level", zaplogger.ExtractErrCtx(errors.Ctx().Loc(2).Just(err))...)
		return
	}

	// create DI & BL
	di := internal.NewDI(logger)
	bli := bl.NewBL(di)

	finished := make(chan bool)
	server := http2.NewHTTPServer(logger, bli, finished)
	server.Run(logger)

	<-finished
}
