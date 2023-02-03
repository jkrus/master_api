package config

import (
	"time"

	tracingCfg "github.com/jkrus/master_api/pkg/tracing/v2/config"
)

func (c *Config) LocalAgentHost() string {
	return c.JaegerLocalAgentHost
}

func (c *Config) LocalAgentPort() string {
	return c.JaegerLocalAgentPort
}

func (c *Config) Type() string {
	return c.JaegerTyp
}

func (c *Config) SamplingServerURL() string {
	return c.JaegerSamplingServerURL
}

func (c *Config) Param() float64 {
	return c.JaegerParam
}

func (c *Config) BufferFlushInterval() time.Duration {
	return c.JaegerBufferFlushInterval
}

func (c *Config) AttemptReconnectInterval() time.Duration {
	return c.JaegerAttemptReconnectInterval
}

func (c *Config) SamplingRefreshInterval() time.Duration {
	return c.JaegerSamplingRefreshInterval
}

func (c *Config) QueueSize() int {
	return c.JaegerQueueSize
}

func (c *Config) LogSpans() bool {
	return c.JaegerLogSpans
}

func (c *Config) Enabled() bool {
	return Get().JaegerEnabled
}

func (c *Config) Get() tracingCfg.TraceConfiger {
	return GetForTracer()
}

func GetForTracer() tracingCfg.TraceConfiger {
	c := globalConfig.CopyData().(Config)
	return &c
}
