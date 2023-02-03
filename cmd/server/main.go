package main

import (
	"context"

	"github.com/jkrus/master_api/internal/config"
	"github.com/jkrus/master_api/pkg/errors"
	"github.com/jkrus/master_api/pkg/tracing/v2"

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
		logger.S("logger creating error", zaplogger.ExtractErrCtx(errors.Ctx().Loc(2).Just(err))...)
		return
	}

	ctx := context.Background()
	// init config
	err = config.Init(ctx, appName, logger, loggerLevel)
	if err != nil {
		logger.Error("config init error", zaplogger.ExtractErrCtx(errors.Ctx().Loc(2).Just(err))...)
		return
	}

	// change log level
	settings := config.Get()
	err = loggerLevel.UnmarshalText([]byte(settings.LogLevel))
	if err != nil {
		logger.Error("can't change log level", zaplogger.ExtractErrCtx(errors.Ctx().Loc(2).Just(err))...)
		return
	}

	// create tracer
	tracer, err := tracing.NewTracer(logger, appName, &settings, nil)
	if err != nil {
		logger.Error("tracer creating error", zaplogger.ExtractErrCtx(errors.Ctx().Loc(2).Just(err))...)
	}
	defer tracer.Close()

	// create DI & BL
	di := internal.NewDI(logger)
	bli := bl.NewBL(di)

	finished := make(chan bool)
	server := http2.NewHTTPServer(logger, bli, finished)
	server.Run(logger)

	<-finished
}
