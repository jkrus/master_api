package postgre

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jkrus/master_api/internal/config"
	db2 "github.com/jkrus/master_api/internal/stores/db"

	gormLogger "gorm.io/gorm/logger"

	_ "github.com/golang-migrate/migrate/source/file" // registration of migrations files
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDBConnection(logger *zap.Logger) (*gorm.DB, error) {
	logger = logger.Named("PostgreSQL_DB")
	settings := config.GetConfig()
	dataSourceName := makeDataSourceName(*settings)
	var (
		db  *gorm.DB
		err error
	)

	if db, err = connectWithGORMRetry(dataSourceName, logger, *settings); err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(3 * time.Minute)

	if !settings.AutoMigrate {
		return db, nil
	}

	if err := db2.ApplyAutoMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func connectWithGORMRetry(connectionString string, logger *zap.Logger, settings config.Config) (*gorm.DB, error) {
	ticker := time.NewTicker(1 * time.Nanosecond)
	timeout := time.After(15 * time.Minute)
	seconds := 1
	try := 0
	for {
		select {
		case <-ticker.C:
			try++
			ticker.Stop()
			client, err := connectWithGORM(connectionString, logger, settings)
			if err != nil {
				logger.Warn(fmt.Sprintf("не удалось установить соединение с PostgreSQL, попытка № %d", try), zap.Error(err))

				ticker = time.NewTicker(time.Duration(seconds) * time.Second)
				seconds *= 2
				if seconds > 60 {
					seconds = 60
				}
				continue
			}

			logger.Debug("соединение с PostgreSQL успешно установлено")
			return client, nil
		case <-timeout:
			return nil, errors.New("PostgreSQL: connection timeout")
		}
	}
}

func makeDataSourceName(settings config.Config) string {
	parameters := map[string]string{
		"host":     settings.DBHost,
		"port":     settings.DBPort,
		"user":     settings.DBUser,
		"password": settings.DBPass,
		"dbname":   settings.DBName,
		"sslmode":  settings.DBSSLMode,
	}

	var pairs []string
	for key, value := range parameters {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(pairs, " ")
}

func connectWithGORM(dataSourceName string, logger *zap.Logger, settings config.Config) (*gorm.DB, error) {
	logLevel := gormLogger.Warn
	if settings.TraceSQLCommands {
		logLevel = gormLogger.Info
	}

	return gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: db2.NewLogger(logger, gormLogger.Config{
			// временной зазор определения медленных запросов SQL
			SlowThreshold: time.Duration(settings.SQLSlowThreshold) * time.Second,
			LogLevel:      logLevel,
			Colorful:      false,
		}),
		AllowGlobalUpdate: true,
	})
}
