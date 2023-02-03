package config

import (
	"context"
	"os"
	"reflect"
	"time"

	"go.uber.org/zap"

	"github.com/jkrus/master_api/pkg/errors"

	"github.com/jkrus/master_api/pkg/zap-logger/v6/fields"

	"github.com/jkrus/master_api/pkg/etcd-configer"
	"github.com/jkrus/master_api/pkg/etcd-configer/configuration"
)

type Config struct {
	// -----------------------------------------------------------------------------
	LogLevel               string        `envconfig:"LOG_LEVEL" default:"info" validate:"oneof=debug info warn error dpanic panic fatal"`
	ServerHost             string        `envconfig:"SERVER_HOST" default:"server"`
	ServerPort             string        `envconfig:"SERVER_PORT" default:"4000"`
	HTTPServerReadTimeOut  time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"10m"`
	HTTPServerWriteTimeOut time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"13m"`
	PayloadSoftLimit       int           `envconfig:"PAYLOAD_SAFE_LIMIT" default:"5120"`    // 5 КиБ
	PayloadHardLimit       int           `envconfig:"PAYLOAD_HARD_LIMIT" default:"5242880"` // 5 МиБ
	PayloadQuantityLimit   int           `envconfig:"PAYLOAD_QUANTITY_LIMIT" default:"10000"`

	// CacheRequestTTL - время жизни кэшируемого запроса
	CacheRequestTTL time.Duration `envconfig:"CACHE_REQUEST_TTL" default:"30m"`
	// CacheResponseTTL - время жизни кэшируемого ответа
	CacheResponseTTL time.Duration `envconfig:"CACHE_RESPONSE_TTL" default:"30m"`

	// --------------------------------------- Jaeger -------------------------------------------------
	// Disabled - включение/отключение трейсинга
	JaegerEnabled bool `envconfig:"JAEGER_ENABLED" default:"false"`
	// LocalAgentHost адрес куда нужно отправлять спаны
	JaegerLocalAgentHost string `envconfig:"JAEGER_HOST" default:"jaeger-agent"`
	// LocalAgentPort - порт куда нужно отправлять спаны
	JaegerLocalAgentPort string `envconfig:"JAEGER_PORT" default:"6831"`
	// Type тип сэмлинга, один из: const, probabilistic, rateLimiting или remote
	// по умолчанию remote
	// Соответствует переменной окружения JAEGER_SAMPLER_TYPE
	JaegerTyp string `envconfig:"JAEGER_SAMPLER_TYPE" default:"const"`
	// SamplingServerURL http адрес сэмплируещего сервера(только для remote)
	// Соответствует переменной окружения JAEGER_SAMPLING_ENDPOINT
	JaegerSamplingServerURL string `envconfig:"JAEGER_SAMPLING_ENDPOINT" default:"http://jaeger-agent:5778/sampling"`
	// Param значение, которое будет передано сэплеру.
	// Валидные значения Param :
	// - для "const" сэмплера, 0 или 1 для не сэмплирование ничего или сэмплировать все
	// - для "probabilistic" сэмплера, значение от 0 до 1
	// - для "rateLimiting" сэмплера, количество спанов в секунду
	// - для "remote" сэмплера - то же самое что и probabilistic, до того момента пока
	//   он не получит значение от удаленного сервера
	// Соответствует переменной окружения JAEGER_SAMPLER_PARAM
	JaegerParam float64 `envconfig:"JAEGER_SAMPLER_PARAM" default:"1"`
	// BufferFlushInterval контролирует как часто буфер со спанами будут принудительно
	// отправлены на сервер, даже если буфер пуст
	JaegerBufferFlushInterval time.Duration `envconfig:"JAEGER_REPORTER_FLUSH_INTERVAL" default:"10s"`
	// AttemptReconnectInterval контролирует с какой частотой клиент проверяет смену адреса сервиса трейсинга.
	JaegerAttemptReconnectInterval time.Duration `envconfig:"JAEGER_REPORTER_ATTEMPT_RECONNECT_INTERVAL" default:"30s"`
	// SamplingRefreshInterval описывает частоту опроса сэплируещего сервера(только для remote)
	// по умолчанию 1 минута
	JaegerSamplingRefreshInterval time.Duration `envconfig:"JAEGER_SAMPLER_REFRESH_INTERVAL" default:"1m"`
	// QueueSize размер очереди спанов - если очередь полна то спаны будут отброшены
	JaegerQueueSize int `envconfig:"JAEGER_REPORTER_MAX_QUEUE_SIZE" default:"100"`
	// LogSpans когда это значение равно истине, в логи будут писаться сообщение
	// о том какие спаны были отправлены на сервер.
	JaegerLogSpans bool `envconfig:"JAEGER_REPORTER_LOG_SPANS" default:"false"`
}

func (config Config) PayloadConfig() fields.PayloadConfig {
	payloadConfig := fields.DefaultPayloadConfig
	payloadConfig.SoftLimit = config.PayloadSoftLimit
	payloadConfig.HardLimit = config.PayloadHardLimit
	payloadConfig.QuantityLimit = config.PayloadQuantityLimit

	return payloadConfig
}

var (
	globalConfig configuration.StructuralConfiguration
)

func Get() Config {
	return globalConfig.CopyData().(Config)
}

func Init(ctx context.Context, appName string, logger *zap.Logger, loggerLevel zap.AtomicLevel) error {
	etcdEndpoint := os.Getenv("ETCD_ENDPOINT")
	if etcdEndpoint == "" {
		etcdEndpoint = "etcd:2379"
	}
	return initCfg(
		ctx,
		etcdEndpoint,
		// options
		etcdconfiger.WithoutStorage(os.Getenv("ETCD_DISABLED") == "TRUE"),
		etcdconfiger.WithLoadingPrefix("root/"+appName+"/"),
		// options/timings
		etcdconfiger.WithDialTimeout(5*time.Second),
		etcdconfiger.WithLoadingTimeout(5*time.Second),
		etcdconfiger.WithRetryDelay(5*time.Second),
		// options/handlers
		etcdconfiger.WithLogger(logger.Named("etcd")),
		etcdconfiger.WithUpdateHandler(getChangeLoggerLevel(logger, loggerLevel)),
	)
}

const (
	ErrUnableToUpdateLogLevel = "не могу сменить log level"
	ErrInitConfig             = "ошибка инициализации конфига"
)

func getChangeLoggerLevel(logger *zap.Logger, loggerLevel zap.AtomicLevel) configuration.UpdateHandler {
	return func(name string, value interface{}) {
		if name == "LOG_LEVEL" {
			if err := loggerLevel.UnmarshalText([]byte(value.(string))); err != nil {
				logger.
					With(fields.AnyWithSuffix("level", value, fields.FunctionNameSuffix, fields.CounterSuffix)).
					Error(ErrUnableToUpdateLogLevel)
			}
		}
	}
}

func initCfg(
	ctx context.Context,
	etcdEndpoint string,
	options ...etcdconfiger.ConfigurationOption,
) error {
	config, err := etcdconfiger.Init(ctx, reflect.TypeOf(Config{}), etcdEndpoint, options...)
	if err != nil {
		return errors.Wrapf(err, ErrInitConfig)
	}
	globalConfig = config
	return nil
}
