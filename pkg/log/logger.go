package log

import (
	"io"
	"sync"
)

type NoopLogger NoopWriter

func (l NoopLogger) Trace() Writer        { return NoopWriter(l) }
func (l NoopLogger) Debug() Writer        { return NoopWriter(l) }
func (l NoopLogger) Info() Writer         { return NoopWriter(l) }
func (l NoopLogger) Warning() Writer      { return NoopWriter(l) }
func (l NoopLogger) Error() Writer        { return NoopWriter(l) }
func (l NoopLogger) V(_ Verbosity) Writer { return NoopWriter(l) }

type StreamLogger struct {
	mu sync.Mutex

	v Verbosity

	sinks []Writer
	all   Writer
}

func (l *StreamLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	l.all = NewStreamWriter(w)
	l.setVerbosity(l.v)
	l.mu.Unlock()
}

func (l *StreamLogger) SetVerbosity(v Verbosity) {
	l.mu.Lock()
	l.setVerbosity(v)
	l.mu.Unlock()
}

func (l *StreamLogger) setVerbosity(v Verbosity) {
	l.v = v
	l.sinks = []Writer{
		NoopStreamWriter, // OFF
		NoopStreamWriter, // TRACE
		NoopStreamWriter, // DEBUG
		NoopStreamWriter, // INFO
		NoopStreamWriter, // WARNING
		NoopStreamWriter, // ERROR
	}

	switch {
	case v >= TRACE:
		l.sinks[TRACE] = l.all
		fallthrough
	case v >= DEBUG:
		l.sinks[DEBUG] = l.all
		fallthrough
	case v >= INFO:
		l.sinks[INFO] = l.all
		fallthrough
	case v >= WARNING:
		l.sinks[WARNING] = l.all
		fallthrough
	case v >= ERROR:
		l.sinks[ERROR] = l.all
	}
}

func (l *StreamLogger) V(v Verbosity) Writer {
	if int(v) < len(l.sinks) {
		return l.sinks[v]
	}

	return l.sinks[TRACE]
}

func (l *StreamLogger) Error() Writer {
	return l.sinks[ERROR]
}

func (l *StreamLogger) Warning() Writer {
	return l.sinks[WARNING]
}

func (l *StreamLogger) Info() Writer {
	return l.sinks[INFO]
}

func (l *StreamLogger) Debug() Writer {
	return l.sinks[DEBUG]
}

func (l *StreamLogger) Trace() Writer {
	return l.sinks[TRACE]
}

func NewStreamLogger(output io.Writer) *StreamLogger {
	result := &StreamLogger{
		v: OFF,
	}

	result.SetOutput(output)

	return result
}

var (
	DisabledNoopLogger = NoopLogger(DisabledNoopWriter)
)
