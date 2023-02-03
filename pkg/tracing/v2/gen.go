package tracing

//go:generate mockgen -source=./config/config.go -destination=./mocks/configger.go -package=mocks TraceConfiger
//go:generate mockgen -destination=./mocks/opentracing.go -package=mocks github.com/opentracing/opentracing-go Tracer,Span

//go:generate gomarkdoc -o ./README.md . ./config
