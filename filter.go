package logger

import "context"

var EmptyFilter = func(context.Context, map[string]string) bool { return false }

var DefaultFilter = NewLevelFilter(LevelInfo)
