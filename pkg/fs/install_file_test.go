package fs

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestDirMode(t *testing.T) {
	type testCase struct {
		name         string
		input        os.FileMode
		assertOutput os.FileMode
	}

	run := func(t *testing.T, tc testCase) {
		got := dirMode(tc.input)

		assert.Equal(t, got.String(), tc.assertOutput.String())
	}
	testCases := []testCase{
		{
			name: "empty",
		},
		{
			name:  "only executable",
			input: 0111,
		},
		{
			name:  "only writable",
			input: 0222,
		},
		{
			name:         "world readable",
			input:        0644,
			assertOutput: 0755,
		},
		{
			name:         "group readable",
			input:        0640,
			assertOutput: 0750,
		},
		{
			name:         "user writable",
			input:        0600,
			assertOutput: 0700,
		},
		{
			name:         "user readable",
			input:        0400,
			assertOutput: 0700,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
