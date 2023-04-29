package cli

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestConfigPathsDirectory(t *testing.T) {
	type testCase struct {
		name         string
		input        []string
		assertOutput string
	}

	run := func(t *testing.T, tc testCase) {
		out := ConfigPaths(tc.input)
		got := out.Directory()

		assert.Equal(t, got, tc.assertOutput)
	}
	testCases := []testCase{
		{
			name:  "empty",
			input: []string{},
		},
		{
			name: "first match",
			input: []string{
				"cmd/edit.go",
				"testdata/resolver.toml",
				"./config.go",
			},
			assertOutput: "cmd",
		},
		{
			name: "last resort",
			input: []string{
				"does/not/exist.cfg",
				"this/one/neither.txt",
				"also/missing/file.mp4",
			},
			assertOutput: "does/not",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
