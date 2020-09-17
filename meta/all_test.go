package meta

import (
	"context"
	"testing"
)

func Test_Logger_Meta_All(t *testing.T) {
	ctx := context.Background()

	{
		m := All(ctx)
		if m == nil {
			t.Fatal("meta must not be nil since it should always be initialized")
		}
		if len(m) != 0 {
			t.Fatal("meta must be empty since it has not been set")
		}
	}

	{
		ctx = Add(ctx, "foo", "bar")
		ctx = Add(ctx, "foo", "changed")
		ctx = Add(ctx, "bar", "test")
	}

	{
		m := All(ctx)
		if len(m) != 2 {
			t.Fatal("meta must have 2 key-value pairs since they have been set")
		}
	}

	// We want to make sure that read access to the logger meta does not allow
	// write access to the underlying data structure. So in case someone can
	// overwrite the internal state of the logger meta by modifying the returned
	// map received by calling All, we would have a problem which we do not want
	// to have.
	{
		m := All(ctx)
		m["hack"] = "bad"
	}

	{
		ctx = Rem(ctx, "foo")
		ctx = Rem(ctx, "bar")
	}

	{
		m := All(ctx)
		if m == nil {
			t.Fatal("meta must not be nil since it should always be initialized")
		}
		if len(m) != 0 {
			t.Fatal("meta must be empty since all key-value pairs have been removed")
		}
	}
}
