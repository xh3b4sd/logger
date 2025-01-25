package logger

import "context"

type fake struct{}

func (r fake) Log(pai ...string) {}

func (r fake) LogCtx(ctx context.Context, pai ...string) {}
