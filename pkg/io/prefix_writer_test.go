package io

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func SkipTestPrefixWriter(t *testing.T) {
	var buf strings.Builder
	out := NewPrefixWriter(&buf, []byte("test | "))

	out.Write([]byte("hello"))
	out.Write([]byte(" "))
	out.Write([]byte("world"))
	assert.Equal(t, buf.String(), "test | hello world")

	out.Write([]byte("\r\nhello "))
	out.Write([]byte("\n\nworld\r"))
	out.Write([]byte("hello \n\nworld"))
	out.Write([]byte("\n"))
	assert.Equal(t, buf.String(), `test | hello world
test | hello 
test | 
test | world
test | hello 
test | 
test | world
test | `)
}

func TestPeek(t *testing.T) {
	type testCase struct {
		name         string
		index        int
		input        string
		assertOutput byte
	}

	run := func(t *testing.T, tc testCase) {
		got := Peek([]byte(tc.input), tc.index)

		assert.Equal(t, got, tc.assertOutput)
	}
	testCases := []testCase{
		{
			name: "empty",
		},
		{
			name:  "oob",
			index: 200,
			input: "test",
		},
		{
			name:  "negative",
			index: -1,
			input: "test",
		},
		{
			name:         "first",
			index:        0,
			input:        "test",
			assertOutput: 'e',
		},
		{
			name:         "last",
			index:        2,
			input:        "test",
			assertOutput: 't',
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
