package logger

import (
	"fmt"
	"runtime"
)

var EmptyCaller = func() string {
	return ""
}

var DefaultCaller = func() string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf("%s:%d", file, line)
}
