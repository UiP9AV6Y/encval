package io

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestReadConfirmation(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		assertOutput bool
		assertError  bool
	}

	run := func(t *testing.T, tc testCase) {
		var dst strings.Builder
		src := strings.NewReader(fmt.Sprintf("%s\n", tc.input))
		got, err := ReadConfirmation(fmt.Sprintf("test case %q", tc.input), src, &dst)

		if tc.assertError {
			assert.Assert(t, err != nil)
		} else {
			assert.NilError(t, err)
			assert.Equal(t, got, tc.assertOutput)
			assert.Equal(t, dst.String(), fmt.Sprintf("test case %q [y/n]: ", tc.input))
		}
	}
	testCases := []testCase{
		{
			name:        "empty",
			assertError: true,
		},
		{
			name:        "invalid",
			input:       "maybe",
			assertError: true,
		},
		{
			name:         "true (long)",
			input:        "yes",
			assertOutput: true,
		},
		{
			name:         "true (short)",
			input:        "y",
			assertOutput: true,
		},
		{
			name:         "false (long)",
			input:        "no",
			assertOutput: false,
		},
		{
			name:         "false (short)",
			input:        "n",
			assertOutput: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
