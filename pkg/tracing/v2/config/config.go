// Package config имплементация конфигера jaeger трейсера
package config

import (
	"net"
	"time"

	"github.com/uber/jaeger-client-go/config"

	"github.com/jkrus/master_api/pkg/errors"
)

// TraceConfiger интерфейс конфигурирования трейсера
type TraceConfiger interface {
	// LocalAgentHost адрес куда нужно отправлять спаны
	// Соответствует переменной JAEGER_AGENT_HOST, по умолчанию localhost
	LocalAgentHost() string
	// LocalAgentPort - порт куда нужно отправлять спаны
	// Соответствует переменной JAEGER_AGENT_PORT, по умолчанию 6831
	LocalAgentPort() string
	// Type тип сэмлинга, один из: const, probabilistic, rateLimiting или remote
	// по умолчанию remote
	// Соответствует переменной JAEGER_SAMPLER_TYPE, по умолчанию remote
	Type() string
	// SamplingServerURL http адрес сэмплируещего сервера(только для remote)
	// Соответствует переменной JAEGER_SAMPLING_ENDPOINT, по умолчанию http://localhost:5778/sampling
	SamplingServerURL() string
	// Param значение которое будет передано сэплеру.
	// Валидные значения Param :
	// - для "const" сэмплера, 0 или 1 для не сэмплирование ничего или сэмплировать все
	// - для "probabilistic" сэмплера, значание от 0 до 1
	// - для "rateLimiting" сэмплера, количество спанов в секунду
	// - для "remote" сэмплера - тоже самое что и probabilistic, до того момента пока
	//   он не получить значение от удаленного сервера
	// Соответствует переменной JAEGER_SAMPLER_PARAM, по умолчанию 1
	Param() float64
	// BufferFlushInterval контролирует как часто буфер со спанами будут принудительно
	// отправлены на сервер, даже если буфер пуст
	// Соответствует переменной JAEGER_REPORTER_FLUSH_INTERVAL, по умолчанию 10s
	BufferFlushInterval() time.Duration
	// AttemptReconnectInterval контролирует к с какой частотой клиент проверяет смену адреса сервиса трейсинга.
	// Соответствует переменной JAEGER_REPORTER_ATTEMPT_RECONNECT_INTERVAL, по умолчанию 30s
	AttemptReconnectInterval() time.Duration
	// SamplingRefreshInterval описывает частоту опроса сэплируещего сервера(только для remote)
	// по умолчанию 1 минута
	// Соответствует переменной JAEGER_SAMPLER_REFRESH_INTERVAL, по умолчанию 1m
	SamplingRefreshInterval() time.Duration
	// QueueSize размер очереди спанов - если очередь полна то спаны будут отброшены
	// Соответствует переменной JAEGER_REPORTER_MAX_QUEUE_SIZE, по умолчанию 100
	QueueSize() int
	// LogSpans когда это значение равно истине, в логи будут писаться сообщение
	// о том какие спаны были отправлены на сервер.
	// Соответствует переменной JAEGER_REPORTER_LOG_SPANS, по умолчанию false
	LogSpans() bool
	// Enabled - включение/отключение трейсинга
	// Соответствует переменной JAEGER_ENABLED, по умолчанию false
	Enabled() bool

	// Get - возвращает свежую конфигурацию джагера
	Get() TraceConfiger
}

// JaegerConfiguration возвращает конфиг jaeger клиента
func JaegerConfiguration(serviceName string, c TraceConfiger) (*config.Configuration, error) {
	if c != nil {
		jcfg := &config.Configuration{
			ServiceName: serviceName,
			Disabled:    !c.Enabled(),
		}

		if jcfg.Disabled {
			return jcfg, nil
		}

		jcfg.Reporter, jcfg.Sampler = &config.ReporterConfig{}, &config.SamplerConfig{}
		jcfg.Reporter.LocalAgentHostPort = net.JoinHostPort(c.LocalAgentHost(), c.LocalAgentPort())
		jcfg.Reporter.BufferFlushInterval = c.BufferFlushInterval()
		jcfg.Reporter.AttemptReconnectInterval = c.AttemptReconnectInterval()
		jcfg.Reporter.LogSpans = c.LogSpans()
		jcfg.Reporter.QueueSize = c.QueueSize()
		jcfg.Sampler.Type = c.Type()
		jcfg.Sampler.Param = c.Param()
		jcfg.Sampler.SamplingRefreshInterval = c.SamplingRefreshInterval()
		jcfg.Sampler.SamplingServerURL = c.SamplingServerURL()

		return jcfg, nil
	}

	jcfg, err := config.FromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "get jaeger configuration from env")
	}

	jcfg.ServiceName = serviceName

	return jcfg, nil
}
