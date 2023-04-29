package crypto

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestDataCopy(t *testing.T) {
	type testCase struct {
		name  string
		input Data
	}

	run := func(t *testing.T, tc testCase) {
		got := tc.input.Copy()

		assert.Assert(t, got != nil)
		assert.Equal(t, len(got), len(tc.input))
	}
	testCases := []testCase{
		{
			name:  "empty",
			input: Data([]byte{}),
		},
		{
			name:  "common usage",
			input: Data([]byte("test")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
