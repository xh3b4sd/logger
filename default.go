package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/go-stack/stack"
)

var DefaultCaller = func() string {
	return fmt.Sprintf("%+v", stack.Caller(2))
}

var DefaultFilter = NewLevelFilter(LevelInfo)

var DefaultFormatter = JSONFormatter

var DefaultTimer = func() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

var DefaultWriter = os.Stdout
