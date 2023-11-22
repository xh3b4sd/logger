package logger

import "context"

type fake struct{}

func (r fake) Log(ctx context.Context, pai ...string) {}
