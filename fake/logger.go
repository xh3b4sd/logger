package fake

import "context"

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Log(ctx context.Context, keyVals ...string) {}
