# logger

Simple structured logging with severity awareness. Logging is one essential
observability technique in order to not fly blindly. Structuring logs is simple
and yet powerful because it supports different formats that can be addapted to
any application's needs.



### Log Level Awareness

There are 4 log levels supported, namely `debug`, `info`, `warning` and `error`.
Logs are emitted based on filter configurations provided to the logger creation.
Below is shown which logger setting causes which behaviour of emitted logs.

logger setting | emitted logs
---|---
`debug` | `debug`, `info`, `warning`, `error`
`info` | `info`, `warning`, `error`
`warning` | `warning`, `error`
`error` | `error`



### Logging In Code

Below is an example of emitting `info` logs providing information of expected
business logic behaviour.

```golang
r.logger.Log(ctx, "level", "info", "message", "something important happened")
```

Below is an example of emitting `error` logs providing information of unexpected
business logic behaviour. Note that stack trace logging is wel integrated with
[tracer](https://github.com/xh3b4sd/tracer).

```golang
r.logger.Log(ctx, "level", "error", "message", "something bad happened", "stack", tracer.JSON(err))
```



### Log Line Printing

Using `logger` prints log lines like shown in the example below. Note that you
can configure the `io.Writer` that actually handles the byte streams. The
default is configured to be `os.Stdout`.

```json
{
	"caller": "--REPLACED--/logger_test.go:129",
	"level": "info",
	"message": "something important happened",
	"time": "2006-01-02 15:04:05"
}
```
