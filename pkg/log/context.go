package log

import (
	"context"
	"errors"
)

var ErrNoLogger = errors.New("no Logger found in Context")

type loggerKey struct{}

// NewContext creates a new context with logger information attached.
func NewContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext returns the logger information in ctx if it exists.
func FromContext(ctx context.Context) (Logger, error) {
	l, ok := ctx.Value(loggerKey{}).(Logger)
	if !ok {
		return nil, ErrNoLogger
	}

	return l, nil
}
