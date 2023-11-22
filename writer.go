package logger

import (
	"io"
	"os"
)

var EmptyWriter = io.Discard

var DefaultWriter = os.Stdout
