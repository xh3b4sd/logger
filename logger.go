package logger

import (
	"context"
	"fmt"
	"io"

	"github.com/xh3b4sd/logger/meta"
	"github.com/xh3b4sd/tracer"
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
	Caller func() string
	Filter func(context.Context, map[string]string) bool
	Format func(context.Context, map[string]string) string
	Timer  func() string
	Writer io.Writer
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
	if config.Format == nil {
		config.Format = DefaultFormatter
	}
	if config.Timer == nil {
		config.Timer = DefaultTimer
	}
	if config.Writer == nil {
		config.Writer = DefaultWriter
	}

	return &Logger{
		cal: config.Caller,
		fil: config.Filter,
		frm: config.Format,
		tim: config.Timer,
		wri: config.Writer,
	}
}

func (l *Logger) Log(pai ...string) {
	if len(pai)%2 != 0 {
		_, err := fmt.Fprintln(l.wri, "Log() received uneven amount of key-value pairs")
		if err != nil {
			tracer.Panic(err)
		}

		return
	}

	// At first we need to have a map full of all the key-value pairs of the
	// current log line.
	all := map[string]string{}
	for i := 1; i < len(pai); i += 2 {
		all[pai[i-1]] = pai[i]
	}

	// We check if the current log line should be emitted or filtered. The
	// configured filter tells us if we should proceed or not.
	if l.fil(nil, all) {
		return
	}

	// Set the additional magic key-value pairs to get e.g. the log caller and
	// time.
	{
		cal := l.cal()
		if cal != "" {
			all[KeyCal] = cal
		}

		tim := l.tim()
		if tim != "" {
			all[KeyTim] = tim
		}
	}

	// What we get from the configured formatter can simply be written to the
	// configured writer.
	{
		_, err := fmt.Fprintln(l.wri, l.frm(nil, all))
		if err != nil {
			tracer.Panic(err)
		}
	}
}

func (l *Logger) LogCtx(ctx context.Context, pai ...string) {
	if len(pai)%2 != 0 {
		_, err := fmt.Fprintln(l.wri, "LogCtx() received uneven amount of key-value pairs")
		if err != nil {
			tracer.Panic(err)
		}

		return
	}

	// At first we need to have a map full of all the key-value pairs of the
	// current log line.
	all := meta.All(ctx)
	for i := 1; i < len(pai); i += 2 {
		all[pai[i-1]] = pai[i]
	}

	// We check if the current log line should be emitted or filtered. The
	// configured filter tells us if we should proceed or not.
	if l.fil(ctx, all) {
		return
	}

	// Set the additional magic key-value pairs to get e.g. the log caller and
	// time.
	{
		cal := l.cal()
		if cal != "" {
			all[KeyCal] = cal
		}

		tim := l.tim()
		if tim != "" {
			all[KeyTim] = tim
		}
	}

	// What we get from the configured formatter can simply be written to the
	// configured writer.
	{
		_, err := fmt.Fprintln(l.wri, l.frm(ctx, all))
		if err != nil {
			tracer.Panic(err)
		}
	}
}
