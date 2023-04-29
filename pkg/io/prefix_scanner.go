package io

import (
	"bufio"
	"bytes"
	libio "io"
)

// ScanLines is a splitter function for a bufio.Scanner instance. it is the
// same implementation as the one from the golang library, except it preserves
// the terminating newline characters
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, NL); i >= 0 {
		if Peek(data, i) == CR {
			return i + 2, data[0 : i+2], nil
		}
		return i + 1, data[0 : i+1], nil
	} else if i := bytes.IndexByte(data, CR); i >= 0 {
		if Peek(data, i) == NL {
			return i + 2, data[0 : i+2], nil
		}
		return i + 1, data[0 : i+1], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

type PrefixScanner struct {
	scanner     *bufio.Scanner
	reader      libio.Reader
	prefix      []byte
	passthrough bool
}

func NewPrefixScanner(reader libio.Reader, prefix []byte) *PrefixScanner {
	scanner := bufio.NewScanner(reader)
	result := &PrefixScanner{
		scanner: scanner,
		reader:  reader,
		prefix:  prefix,
	}

	scanner.Split(ScanLines)

	return result
}

func (r *PrefixScanner) Next() ([]byte, error) {
	if !r.passthrough {
		for r.scanner.Scan() {
			line := r.scanner.Bytes()
			if !bytes.HasPrefix(line, r.prefix) {
				r.passthrough = true
				return line, r.scanner.Err()
			}
		}
	}

	if !r.scanner.Scan() {
		return nil, r.scanner.Err()
	}

	return r.scanner.Bytes(), r.scanner.Err()
}
