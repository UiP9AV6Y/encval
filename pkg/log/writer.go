package log

import (
	"fmt"
	"io"
)

type NoopWriter bool

func (w NoopWriter) Enabled() bool {
	return bool(w)
}
func (_ NoopWriter) Print(_ ...interface{})            {}
func (_ NoopWriter) Printf(_ string, _ ...interface{}) {}
func (_ NoopWriter) Println(_ ...interface{})          {}

type StreamWriter struct {
	io.Writer
}

func NewDisabledStreamWriter() *StreamWriter {
	result := &StreamWriter{
		Writer: io.Discard,
	}

	return result
}

func NewStreamWriter(w io.Writer) *StreamWriter {
	result := &StreamWriter{
		Writer: w,
	}

	return result
}

func (w *StreamWriter) Enabled() bool {
	return w.Writer != io.Discard
}

func (w *StreamWriter) Print(a ...interface{}) {
	fmt.Fprint(w, a...)
}

func (w *StreamWriter) Printf(format string, a ...interface{}) {
	fmt.Fprintf(w, format, a...)
}

func (w *StreamWriter) Println(a ...interface{}) {
	fmt.Fprintln(w, a...)
}

var (
	NoopStreamWriter   = NewDisabledStreamWriter()
	DisabledNoopWriter = NoopWriter(false)
)
