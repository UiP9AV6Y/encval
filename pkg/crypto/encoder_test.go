package crypto

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestEncoding(t *testing.T) {
	type testCase struct {
		name        string
		input       string
		assertError bool
	}

	run := func(t *testing.T, tc testCase) {
		data, err := Encode([]byte(tc.input))
		if tc.assertError {
			assert.Assert(t, err != nil)
			return
		}

		got, err := Decode(data)
		assert.NilError(t, err)
		assert.Equal(t, string(got), tc.input)
	}
	testCases := []testCase{
		{
			name: "empty",
		},
		{
			name:  "space",
			input: " ",
		},
		{
			name:  "test",
			input: "Lorem ipsum\ndolor sit amet\n\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
