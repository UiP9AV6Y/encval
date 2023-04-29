package log

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestVerbosityCompare(t *testing.T) {
	type testCase struct {
		name         string
		input        Verbosity
		other        Verbosity
		assertOutput int
	}

	run := func(t *testing.T, tc testCase) {
		got := tc.input.Compare(tc.other)

		assert.Equal(t, got, tc.assertOutput)
	}
	testCases := []testCase{
		{
			name:         "equal",
			input:        OFF,
			other:        OFF,
			assertOutput: 0,
		},
		{
			name:         "more",
			input:        ERROR,
			other:        TRACE,
			assertOutput: -1,
		},
		{
			name:         "less",
			input:        DEBUG,
			other:        INFO,
			assertOutput: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
