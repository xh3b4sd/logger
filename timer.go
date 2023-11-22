package logger

import "time"

var DefaultTimer = func() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}
