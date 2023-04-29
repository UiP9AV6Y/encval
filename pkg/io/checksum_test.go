package io

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCalculateChecksum(t *testing.T) {
	input := "test"
	assertOutput := uint32(3632233996)
	src := strings.NewReader(input)
	got, err := CalculateChecksum(src)

	assert.NilError(t, err)
	assert.Equal(t, got, assertOutput)

	got, err = CalculateChecksum(src)

	assert.NilError(t, err)
	assert.Equal(t, got, uint32(0)) // EOF
}
