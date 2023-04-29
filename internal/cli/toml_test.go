package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"gotest.tools/v3/assert"
)

func TestTOML(t *testing.T) {
	r, err := os.Open("testdata/resolver.toml")
	assert.NilError(t, err)
	assert.Assert(t, r != nil)

	out, err := TOML(r)
	assert.NilError(t, err)
	assert.Assert(t, out != nil)

	type testCase struct {
		name        string
		input       string
		assertValue interface{}
		assertError bool
	}

	run := func(t *testing.T, tc testCase) {
		flag := &kong.Flag{}
		value := &kong.Value{
			Name: tc.input,
			Flag: flag,
		}
		flag.Value = value

		got, err := out.Resolve(nil, nil, flag)

		if tc.assertError {
			assert.Assert(t, err != nil)
		} else {
			assert.NilError(t, err)
			assert.Equal(t, got, tc.assertValue)
		}
	}
	testCases := []testCase{
		{
			name:        "snake case",
			input:       "snake-case",
			assertValue: "test",
		},
		{
			name:        "string value",
			input:       "str",
			assertValue: "value",
		},
		{
			name:        "integer value",
			input:       "int",
			assertValue: int64(1),
		},
		{
			name:        "missing value",
			input:       "nil",
			assertValue: nil,
		},
		{
			name:        "snake case (nested)",
			input:       "nested.snake-case",
			assertValue: "test",
		},
		{
			name:        "string value (nested)",
			input:       "nested.str",
			assertValue: "value",
		},
		{
			name:        "integer value (nested)",
			input:       "nested.int",
			assertValue: int64(1),
		},
		{
			name:        "missing value (nested)",
			input:       "nested.nil",
			assertValue: nil,
		},
		{
			name:        "nestes non-map",
			input:       "ary.0",
			assertValue: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestTOMLDecode(t *testing.T) {
	r := strings.NewReader("invalid == value")

	_, err := TOML(r)
	assert.Assert(t, err != nil)
}
