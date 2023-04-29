package fs

import (
	"os"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestDirHierarchy(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		assertOutput []string
	}

	pwd, err := os.Getwd()
	assert.NilError(t, err)

	run := func(t *testing.T, tc testCase) {
		var want []string
		got := DirHierarchy(tc.input)

		if tc.assertOutput == nil {
			want = strings.Split(pwd, Separator)
		} else {
			want = tc.assertOutput
		}

		assert.Equal(t, len(got), len(want))

		if tc.assertOutput != nil {
			for i := 0; i < len(want); i++ {
				assert.Equal(t, got[i], want[i])
			}
		}
	}
	testCases := []testCase{
		{
			name:  "empty",
			input: "",
		},
		{
			name:  "dot",
			input: ".",
		},
		{
			name:  "root",
			input: "/",
			assertOutput: []string{
				"/",
			},
		},
		{
			name:  "relative",
			input: "/opt/../mnt/./usr/local/bin/",
			assertOutput: []string{
				"/mnt/usr/local/bin",
				"/mnt/usr/local",
				"/mnt/usr",
				"/mnt",
				"/",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
