package io

import (
	"bufio"
	"bytes"
	libio "io"
)

// SeekPrefixedLinesEnd reads lines from the provided reader and adds up
// their length if they are prefixed. The resulting byte count is the
// seek position from the start to skip the affected lines. Lines not starting
// with the prefix stop the search. Any error is from the underlying reader.
func SeekPrefixedLinesEnd(r libio.Reader, prefix []byte) (int64, error) {
	if len(prefix) == 0 {
		return 0, nil
	}

	var result int64
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanLines)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !bytes.HasPrefix(line, prefix) {
			break
		}

		result = result + int64(len(line))
	}

	return result, scanner.Err()
}
