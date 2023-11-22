package logger

import (
	"context"
	"fmt"
	"io"

	"github.com/xh3b4sd/logger/meta"
)

const (
	KeyCal = "caller"
	KeyCod = "code"
	KeyDes = "description"
	KeyDoc = "docs"
	KeyLev = "level"
	KeyMes = "message"
	KeySta = "stack"
	KeyTim = "time"
)

type Config struct {
	Caller    func() string
	Filter    func(context.Context, map[string]string) bool
	Formatter func(context.Context, map[string]string) string
	Timer     func() string
	Writer    io.Writer
}

type Logger struct {
	cal func() string
	fil func(context.Context, map[string]string) bool
	frm func(context.Context, map[string]string) string
	tim func() string
	wri io.Writer
}

func New(config Config) *Logger {
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
		cal: config.Caller,
		fil: config.Filter,
		frm: config.Formatter,
		tim: config.Timer,
		wri: config.Writer,
	}

	return l
}

func (l *Logger) Log(c context.Context, pai ...string) {
	{
		if len(pai)%2 != 0 {
			_, err := fmt.Fprintln(l.wri, "key-value pairs must be complete")
			if err != nil {
				panic(err)
			}
			return
		}
	}

	// At first we need to have a map full of all the key-value pairs of the
	// current log line.
	m := meta.All(c)
	{
		for i := 1; i < len(pai); i += 2 {
			k := pai[i-1]
			v := pai[i]

			m[k] = v
		}
	}

	// We check if the current log line should be emitted or filtered. The
	// configured filter tells us if we should proceed or not.
	{
		if l.fil(c, m) {
			return
		}
	}

	// Set the additional magic key-value pairs to get e.g. the log caller and
	// time.
	{
		cal := l.cal()
		if cal != "" {
			m[KeyCal] = cal
		}

		tim := l.tim()
		if tim != "" {
			m[KeyTim] = tim
		}
	}

	// What we get from the configured formatter can simply be written to the
	// configured writer.
	_, err := fmt.Fprintln(l.wri, l.frm(c, m))
	if err != nil {
		panic(err)
	}
}
