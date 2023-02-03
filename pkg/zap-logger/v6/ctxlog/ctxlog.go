package ctxlog

import (
	"context"

	"go.uber.org/zap"
)

// Use of built-in types for context keys isn't recommended
// to avoid collisions between packages using these.
//
// See for details: https://golang.org/pkg/context/#WithValue
//
type loggerCtxKey struct{}

// GetFromCtx достаёт из контекста логгер
func GetFromCtx(ctx context.Context) *zap.Logger {
	return ctx.Value(loggerCtxKey{}).(*zap.Logger)
}

// AddToCtx добавляет в контекст логгер
func AddToCtx(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

// WithFields ...
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return AddToCtx(ctx, GetFromCtx(ctx).With(fields...))
}
