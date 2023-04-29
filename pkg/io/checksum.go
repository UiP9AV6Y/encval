package io

import (
	"hash/crc32"
	libio "io"
)

func CalculateChecksum(r libio.Reader) (uint32, error) {
	h := crc32.NewIEEE()
	if _, err := libio.Copy(h, r); err != nil {
		return 0, err
	}

	return h.Sum32(), nil
}
