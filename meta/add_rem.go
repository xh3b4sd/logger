package meta

import (
	"context"
)

func Add(ctx context.Context, key, value string) context.Context {
	m, ok := fromContext(ctx)
	if !ok {
		m = new()
		ctx = newContext(ctx, m)
	}

	m.KeyVals[key] = value

	return ctx
}

func Rem(ctx context.Context, key string) context.Context {
	m, ok := fromContext(ctx)
	if ok {
		delete(m.KeyVals, key)
		return ctx
	}

	return ctx
}
