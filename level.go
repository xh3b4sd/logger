package logger

import "fmt"

const (
	LevelKey = "level"
)

const (
	LevelDebug   = "debug"
	LevelInfo    = "info"
	LevelWarning = "warning"
	LevelError   = "error"
)

var (
	levelMapping = map[string]int{
		LevelDebug:   1,
		LevelInfo:    2,
		LevelWarning: 3,
		LevelError:   4,
	}
)

func NewLevelFilter(l string) func(m map[string]string) bool {
	// We compute the filter severity in order to have the basis on which we
	// compare the custom severity given with the current log line. If the given
	// log level does not exist we fail immediately.
	var fs int
	var ok bool
	{
		fs, ok = levelMapping[l]
		if !ok {
			panic(fmt.Sprintf("unknown log level %#q", l))
		}
	}

	return func(m map[string]string) bool {
		// We lookup the custom severity of the log level which we got with the
		// current log line. We check the custom value of the level key we
		// received. It should be e.g. inf or err. If there is no level key
		// given, or if there is no valid log level associated with the given
		// level key, we filter the current log line. This means it will not be
		// emitted.
		var cs int
		{
			cv, ok := m[LevelKey]
			if !ok {
				return true
			}

			cs, ok = levelMapping[cv]
			if !ok {
				return true
			}
		}

		return fs > cs
	}
}
