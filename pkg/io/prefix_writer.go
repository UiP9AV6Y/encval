package io

import (
	libio "io"
)

const (
	// Newline character
	NL = '\n'
	// Carriage Return character
	CR = '\r'
)

type PrefixWriter struct {
	writer libio.Writer
	prefix []byte
	nl, cr bool
}

func NewPrefixWriter(writer libio.Writer, prefix []byte) *PrefixWriter {
	result := &PrefixWriter{
		writer: writer,
		prefix: prefix,
		nl:     true,
		cr:     true,
	}

	return result
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	var i, j int
	for ; i < len(p); i++ {
		if w.nl || w.cr {
			if i > 0 {
				n, err = w.writer.Write(p[j:i])
				if err != nil {
					return
				}

				j = i
			}

			n, err = w.writer.Write(w.prefix)
			if err != nil {
				if n > len(p) {
					n = len(p)
				}
				return
			}

			w.nl = false
			w.cr = false
		}

		if p[i] == NL {
			if Peek(p, i) != CR {
				w.nl = true
			}
		} else if p[i] == CR {
			if Peek(p, i) != NL {
				w.cr = true
			}
		}
	}

	return len(p), nil
}

func Peek(p []byte, n int) byte {
	if n >= 0 && n < len(p)-1 {
		return p[n+1]
	}

	return 0
}
