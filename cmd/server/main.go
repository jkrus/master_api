package main

import (
	"github.com/jkrus/master_api/internal/config"
	"github.com/jkrus/master_api/internal/stores/db"
	"github.com/jkrus/master_api/internal/stores/db/postgre"
	"github.com/jkrus/master_api/internal/stores/minio"
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

	// MinIO
	minioClient, err := minio.NewClient(settings)
	if err != nil {
		logger.Error("can't create MinIO client", zaplogger.ExtractErrCtx(errors.Ctx().Loc(2).Just(err))...)
		return
	}
	minioRepo := minio.NewMinioRepo(logger, minioClient)

	// Postgres DB
	dbClient, err := postgre.OpenDBConnection(logger)
	if err != nil {
		logger.Error("can't create Postgres client", zaplogger.ExtractErrCtx(errors.Ctx().Loc(2).Just(err))...)
		return
	}
	dbRepo := db.NewDBRepo(dbClient)

	// create DI & BL
	di := internal.NewDI(logger, minioRepo, dbRepo)
	bli := bl.NewBL(di)

	finished := make(chan bool)
	server := http2.NewHTTPServer(logger, bli, finished)
	server.Run(logger)

	<-finished
}
