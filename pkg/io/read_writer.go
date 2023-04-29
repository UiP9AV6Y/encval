package io

import (
	libio "io"
)

// ReadWriter implements the io.ReadWriter interface.
// in comparison to bufio.ReadWriter, this implementation
// is completly unbuffered; all calls are passed on directly.
type ReadWriter struct {
	libio.Reader
	libio.Writer
}

// NewReadWriter create as new io.ReadWriter instance using the
// given read and writer instance.
func NewReadWriter(r libio.Reader, w libio.Writer) libio.ReadWriter {
	result := &ReadWriter{
		Reader: r,
		Writer: w,
	}

	return result
}
