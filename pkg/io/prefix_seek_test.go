package io

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSeekPrefixedLinesEnd(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		prefix       string
		assertOutput int64
	}

	run := func(t *testing.T, tc testCase) {
		src := strings.NewReader(tc.input)
		got, err := SeekPrefixedLinesEnd(src, []byte(tc.prefix))

		assert.NilError(t, err)
		assert.Equal(t, got, tc.assertOutput)
	}
	testCases := []testCase{
		{
			name: "empty",
		},
		{
			name:  "empty prefix",
			input: "// hello\nworld\n",
		},
		{
			name:         "no prefix",
			input:        "// hello\nworld\n",
			prefix:       "//",
			assertOutput: 9,
		},
		{
			name:         "all prefix",
			input:        "// hello\n// world\n",
			prefix:       "//",
			assertOutput: 18,
		},
		{
			name:         "no EOL",
			input:        "// hello\n// world",
			prefix:       "//",
			assertOutput: 17,
		},
		{
			name:         "interrupt",
			input:        "// hello\nworld\n// hello\n",
			prefix:       "//",
			assertOutput: 9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
