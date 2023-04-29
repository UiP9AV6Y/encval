package log

import (
	"io"
)

type Verbosity uint8

type Writer interface {
	Enabled() bool
	Print(a ...interface{})
	Printf(format string, a ...interface{})
	Println(a ...interface{})
}

type Logger interface {
	Trace() Writer
	Debug() Writer
	Info() Writer
	Warning() Writer
	Error() Writer
	V(Verbosity) Writer
}

type Controller interface {
	SetOutput(io.Writer)
	SetVerbosity(Verbosity)
}

type LoggerController interface {
	Logger
	Controller
}
