package edit

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func sliceEqual(got, want []string) cmp.Comparison {
	return func() cmp.Result {
		if len(got) != len(want) {
			return cmp.ResultFailure(
				fmt.Sprintf("length of %q != %d", got, len(want)),
			)
		}

		for i := 0; i < len(want); i++ {
			if got[i] != want[i] {
				return cmp.ResultFailure(
					fmt.Sprintf("%q != %q", got[i], want[i]),
				)
			}
		}

		return cmp.ResultSuccess
	}
}

func TestDecorate(t *testing.T) {
	type testCase struct {
		name         string
		format       string
		input        []string
		assertOutput []string
	}

	run := func(t *testing.T, tc testCase) {
		got := decorate(tc.format, tc.input)

		assert.Assert(t, got != nil)
		assert.Assert(t, sliceEqual(got, tc.assertOutput))

	}
	testCases := []testCase{
		{
			name:         "empty",
			input:        []string{},
			assertOutput: []string{},
		},
		{
			name:         "no format",
			input:        []string{"test"},
			assertOutput: []string{"test"},
		},
		{
			name:         "token in input",
			input:        []string{"test %s"},
			assertOutput: []string{"test %s"},
		},
		{
			name:         "enclose",
			format:       "--%s--",
			input:        []string{"unit", "test"},
			assertOutput: []string{"--unit--", "--test--"},
		},
		{
			name:         "no token",
			format:       "--",
			input:        []string{"unit", "test"},
			assertOutput: []string{"--%!(EXTRA string=unit)", "--%!(EXTRA string=test)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestJoin(t *testing.T) {
	type testCase struct {
		name         string
		sep          string
		input        []string
		assertOutput string
	}

	run := func(t *testing.T, tc testCase) {
		got := join(tc.sep, tc.input)

		assert.Equal(t, got, tc.assertOutput)

	}
	testCases := []testCase{
		{
			name:         "empty",
			input:        []string{},
			assertOutput: "",
		},
		{
			name:         "no sep",
			input:        []string{"unit", "test"},
			assertOutput: "unittest",
		},
		{
			name:         "space sep",
			sep:          " ",
			input:        []string{"unit", "test"},
			assertOutput: "unit test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
