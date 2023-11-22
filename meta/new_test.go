package meta

import (
	"context"
	"testing"
)

func Test_Logger_Meta_New(t *testing.T) {
	ctx := context.Background()

	{
		ctx = Add(ctx, "foo", "bar")
	}

	var oth context.Context
	{
		oth = Add(New(ctx), "dif", "tru")
	}

	{
		has := Has(ctx, "foo")
		if !has {
			t.Fatal("foo must exist since it has been set")
		}
	}

	{
		has := Has(ctx, "dif")
		if has {
			t.Fatal("dif must not exist since it has been set in the other context")
		}
	}

	{
		m := All(ctx)
		if m == nil {
			t.Fatal("meta must not be nil since it should always be initialized")
		}
		if len(m) != 1 {
			t.Fatal("meta must be empty since all key-value pairs have been removed")
		}
	}

	{
		has := Has(oth, "foo")
		if !has {
			t.Fatal("foo must exist since it is inherited from the other context")
		}
	}

	{
		has := Has(oth, "dif")
		if !has {
			t.Fatal("dif must exist since it has been set for the other context")
		}
	}

	{
		m := All(oth)
		if m == nil {
			t.Fatal("meta must not be nil since it should always be initialized")
		}
		if len(m) != 2 {
			t.Fatal("meta must be empty since all key-value pairs have been removed")
		}
	}
}
