package logger

import "context"

// Interface implementations emit messages to gather certain runtime
// information.
type Interface interface {
	// Log prints the given key-value pairs.
	Log(pai ...string)
	// LogCtx prints the given key-value pairs. Additionally the given context may
	// provide information injected by the meta package.
	LogCtx(ctx context.Context, pai ...string)
}
