package meta

import (
	"context"
)

func New(ctx context.Context) context.Context {
	n := context.Background()

	for k, v := range All(ctx) {
		n = Add(n, k, v)
	}

	return n
}
