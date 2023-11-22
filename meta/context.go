package meta

import (
	"context"
)

type key string

var loggerMetaKey key = "logger/meta"

type meta struct {
	KeyVals map[string]string
}

func new() *meta {
	return &meta{
		KeyVals: map[string]string{},
	}
}

func newContext(ctx context.Context, v *meta) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, loggerMetaKey, v)
}

// FromContext returns the logger struct, if any.
func fromContext(ctx context.Context) (*meta, bool) {
	v, ok := ctx.Value(loggerMetaKey).(*meta)
	return v, ok
}
