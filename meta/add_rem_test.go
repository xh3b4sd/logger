package meta

import (
	"context"
	"testing"
)

func Test_Logger_Meta_Add_Rem(t *testing.T) {
	ctx := context.Background()

	{
		ok := Has(ctx, "foo")
		if ok {
			t.Fatal("foo must not exist since it has not been set")
		}
	}

	{
		ctx = Add(ctx, "foo", "bar")
	}

	{
		ok := Has(ctx, "foo")
		if !ok {
			t.Fatal("foo must exist since it has been set")
		}
	}

	{
		ctx = Add(ctx, "foo", "changed")
	}

	{
		ok := Has(ctx, "foo")
		if !ok {
			t.Fatal("foo must exist since it has been set")
		}
	}

	{
		ctx = Rem(ctx, "foo")
	}

	{
		ok := Has(ctx, "foo")
		if ok {
			t.Fatal("foo must not exist since it has been removed")
		}
	}
}
