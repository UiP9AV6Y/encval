package edit

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestPreambleWrite(t *testing.T) {
	type testCase struct {
		name string

		*PreambleData
	}

	run := func(t *testing.T, tc testCase) {
		buf := &strings.Builder{}
		err := tc.PreambleData.Write(buf)
		want := fmt.Sprintf("preamble_%s.txt", tc.name)

		assert.NilError(t, err)
		golden.Assert(t, buf.String(), want)

	}
	testCases := []testCase{
		{
			name: "pkcs7",
			PreambleData: &PreambleData{
				AppName:   "encval",
				Prefix:    Prefix,
				Providers: []string{"PKCS7"},
			},
		},
		{
			name: "three",
			PreambleData: &PreambleData{
				AppName:   "test",
				Prefix:    "//",
				Providers: []string{"TEST", "NOOP", "UNIT"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestPreambleLength(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		assertOutput int64
	}

	run := func(t *testing.T, tc testCase) {
		out := strings.NewReader(tc.input)
		got, err := PreambleLength(out)

		assert.NilError(t, err)
		assert.Equal(t, got, tc.assertOutput)
	}
	testCases := []testCase{
		{
			name: "empty",
		},
		{
			name:         "no prefix",
			input:        "hello\nworld\rspec\r\ntest\n\rcase",
			assertOutput: 0,
		},
		{
			name:         "only prefix",
			input:        "# | hello",
			assertOutput: 9,
		},
		{
			name: "1:1",
			input: `# | hello
world`,
			assertOutput: 10,
		},
		{
			name: "intercept",
			input: `# | hello
world
# | hello
world`,
			assertOutput: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
