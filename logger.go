package logger

import (
	"context"
	"fmt"
	"io"

	"github.com/xh3b4sd/logger/meta"
)

const (
	KeyCaller = "caller"
	KeyStack  = "stack"
	KeyTime   = "time"
)

type Config struct {
	Caller    func() string
	Filter    func(m map[string]string) bool
	Formatter func(m map[string]string) string
	Timer     func() string
	Writer    io.Writer
}

type Logger struct {
	caller    func() string
	filter    func(m map[string]string) bool
	formatter func(m map[string]string) string
	timer     func() string
	writer    io.Writer
}

func New(config Config) (*Logger, error) {
	if config.Caller == nil {
		config.Caller = DefaultCaller
	}
	if config.Filter == nil {
		config.Filter = DefaultFilter
	}
	if config.Formatter == nil {
		config.Formatter = DefaultFormatter
	}
	if config.Timer == nil {
		config.Timer = DefaultTimer
	}
	if config.Writer == nil {
		config.Writer = DefaultWriter
	}

	l := &Logger{
		caller:    config.Caller,
		filter:    config.Filter,
		formatter: config.Formatter,
		timer:     config.Timer,
		writer:    config.Writer,
	}

	return l, nil
}

func (l *Logger) Log(ctx context.Context, kvs ...string) {
	if len(kvs)%2 != 0 {
		_, err := fmt.Fprintln(l.writer, "given key-value pairs must be complete for logging")
		if err != nil {
			panic(err)
		}
		return
	}

	// At first we need to have a map full of all the key-value pairs of the
	// current log line.
	m := meta.All(ctx)
	for i := 1; i < len(kvs); i += 2 {
		k := kvs[i-1]
		v := kvs[i]

		m[k] = v
	}

	// We check if the current log line should be emitted or filtered. The
	// configured filter tells us if we should proceed or not.
	if l.filter(m) {
		return
	}

	// Set the additional magic key-value pairs to get e.g. the log caller and
	// time.
	m[KeyCaller] = l.caller()
	m[KeyTime] = l.timer()

	// What we get from the configured formatter can simply be written to the
	// configured writer.
	_, err := fmt.Fprintln(l.writer, l.formatter(m))
	if err != nil {
		panic(err)
	}
}
