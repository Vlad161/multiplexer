package logger

import "context"

var contextKey = struct{}{}

func WithContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, contextKey, log)
}

func FromContext(ctx context.Context) Logger {
	val := ctx.Value(contextKey)
	if log, ok := val.(Logger); ok {
		return log
	}
	return nil
}
