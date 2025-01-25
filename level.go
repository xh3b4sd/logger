package logger

import (
	"context"
	"fmt"

	"github.com/xh3b4sd/tracer"
)

const (
	LevelDebug   = "debug"
	LevelInfo    = "info"
	LevelWarning = "warning"
	LevelError   = "error"
)

var (
	levels = map[string]int{
		LevelDebug:   1,
		LevelInfo:    2,
		LevelWarning: 3,
		LevelError:   4,
	}
)

func NewLevelFilter(l string) func(context.Context, map[string]string) bool {
	// We compute the filter severity in order to have the basis on which we
	// compare the custom severity given with the current log line. If the given
	// log level does not exist we fail immediately.
	var fil int
	var exi bool
	{
		fil, exi = levels[l]
		if !exi {
			tracer.Panic(fmt.Errorf("unknown log level %#q", l))
		}
	}

	return func(_ context.Context, all map[string]string) bool {
		var exi bool

		// We lookup the custom severity of the log level which we got with the
		// current log line. We check the custom value of the level key we
		// received. It should be e.g. info or error. If there is no level key
		// given, or if there is no valid log level associated with the given
		// level key, we filter the current log line. This means it will not be
		// emitted.
		var lev string
		{
			lev, exi = all[KeyLev]
			if !exi {
				return true
			}
		}

		var rnk int
		{
			rnk, exi = levels[lev]
			if !exi {
				return true
			}
		}

		return fil > rnk
	}
}
