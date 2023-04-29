package io

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPrefixScanner(t *testing.T) {
	preamble := "// Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod\n" +
		"// tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,\r" +
		"// quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo\n\r" +
		"// consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse\r\n" +
		"// cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat\n" +
		"// non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n"
	content := "newline\n" +
		"carriage return\r" +
		"nl + cr\n\r" +
		"cr + nl\r\n" +
		"triple nl\n\n" +
		"\nsplit for testing purpose\n" +
		"// keep this\n" +
		"// also this"
	src := strings.NewReader(preamble + content)
	out := NewPrefixScanner(src, []byte("//"))
	var buf strings.Builder

	for {
		got, err := out.Next()

		assert.NilError(t, err)

		if got == nil {
			break
		}

		buf.Write(got)
	}

	assert.Equal(t, buf.String(), content)
}

func TestScanLines(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		atEOF        bool
		assertAdv    int
		assertOutput []byte
		assertError  bool
	}

	run := func(t *testing.T, tc testCase) {
		adv, got, err := ScanLines([]byte(tc.input), tc.atEOF)

		if tc.assertError {
			assert.Assert(t, err != nil)
		} else {
			assert.NilError(t, err)
			assert.Equal(t, adv, tc.assertAdv)

			if tc.assertOutput == nil {
				assert.Assert(t, got == nil)
			} else {
				assert.Equal(t, string(got), string(tc.assertOutput))
			}
		}
	}
	testCases := []testCase{
		{
			name: "empty",
		},
		{
			name:  "empty (EOF)",
			atEOF: true,
		},
		{
			name:         "newline",
			input:        "hello\nworld",
			assertAdv:    6,
			assertOutput: []byte("hello\n"),
		},
		{
			name:         "newlines",
			input:        "hello\n\nworld",
			assertAdv:    6,
			assertOutput: []byte("hello\n"),
		},
		{
			name:         "carriage return",
			input:        "hello\rworld",
			assertAdv:    6,
			assertOutput: []byte("hello\r"),
		},
		{
			name:         "carriage returns",
			input:        "hello\r\rworld",
			assertAdv:    6,
			assertOutput: []byte("hello\r"),
		},
		{
			name:         "CR/NL",
			input:        "hello\r\nworld",
			assertAdv:    7,
			assertOutput: []byte("hello\r\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
