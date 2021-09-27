package ctxpropagation

import (
	"context"

	"go.uber.org/zap"
)

type contextKey int

const (
	loggerContextKey contextKey = iota + 1
	requestIDContextKey
)

func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	logger := ctx.Value(loggerContextKey).(*zap.Logger)
	return logger
}

func SetLoggerForContext(ctx context.Context, logger *zap.Logger) context.Context {
	ctx = context.WithValue(ctx, loggerContextKey, logger)
	return ctx
}

func GetRequestIDFromContext(ctx context.Context) string {
	requestID := ctx.Value(loggerContextKey).(string)
	return requestID
}

func SetRequestIDForContext(ctx context.Context, requestID string) context.Context {
	ctx = context.WithValue(ctx, loggerContextKey, requestID)
	return ctx
}
