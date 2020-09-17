package meta

import (
	"context"
)

func Has(ctx context.Context, key string) bool {
	m, ok := fromContext(ctx)
	if !ok {
		return false
	}

	_, has := m.KeyVals[key]

	return has
}
