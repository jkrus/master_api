package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	jzap "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"

	tracingcfg "github.com/jkrus/master_api/pkg/tracing/v2/config"
)

type ITracer interface {
	Get() opentracing.Tracer
	Close()
}

type TracerWrap struct {
	tracer    opentracing.Tracer
	closeFunc CloseFunc
}

// CloseFunc метод закрытия трейсера
type CloseFunc func()

// FromConfig конструктор opentracing.Tracer из конфига
func FromConfig(logger *zap.Logger, serviceName string, tConf tracingcfg.TraceConfiger, opts []config.Option) (opentracing.Tracer, CloseFunc, error) {
	conf, err := tracingcfg.JaegerConfiguration(serviceName, tConf)
	if err != nil {
		return nil, nil, err
	}

	opts = append(opts, config.Logger(jzap.NewLogger(logger)))

	tracer, closer, err := conf.NewTracer(opts...)
	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, func() {
		err = closer.Close()
		if err != nil {
			logger.Error("tracer close error", zap.Error(err))
		}
	}, nil
}

func NewTracer(logger *zap.Logger, serviceName string, tConf tracingcfg.TraceConfiger, opts []config.Option) (ITracer, error) {
	tracer, closeFunc, err := FromConfig(logger, serviceName, tConf, opts)
	if err != nil {
		return nil, err
	}
	return &TracerWrap{
		tracer:    tracer,
		closeFunc: closeFunc,
	}, nil
}

func (tw TracerWrap) Get() opentracing.Tracer {
	return tw.tracer
}

func (tw TracerWrap) Close() {
	tw.closeFunc()
}
