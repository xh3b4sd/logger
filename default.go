package logger

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/go-stack/stack"
)

var DefaultCaller = func() string {
	return fmt.Sprintf("%+v", stack.Caller(2))
}

var DefaultFilter = NewLevelFilter(LevelInfo)

var DefaultFormatter = func(m map[string]string) string {
	// Once we have the full map of key-value pairs we need the keys only in
	// order to sort them. We want the emitted log line to be alphabetically
	// ordered by keys.
	var keys []string
	{
		for k, _ := range m {
			keys = append(keys, k)
		}

		sort.Sort(sort.StringSlice(keys))
	}

	// Below we compute a JSON string for the structured log line we want to
	// emit.
	var s string
	{
		s += "{ "

		for i, k := range keys {
			s += "\""
			s += k
			s += "\""
			s += ":"
			s += "\""
			s += m[k]
			s += "\""

			if i+1 < len(keys) {
				s += ", "
			}
		}

		s += " }"
	}

	return s
}

var DefaultTimer = func() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

var DefaultWriter = os.Stdout
