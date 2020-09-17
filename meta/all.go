package meta

import (
	"context"
)

func All(ctx context.Context) map[string]string {
	n := map[string]string{}

	m, ok := fromContext(ctx)
	if !ok {
		return n
	}

	for k, v := range m.KeyVals {
		n[k] = v
	}

	return n
}
