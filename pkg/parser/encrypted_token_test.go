package parser

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

func TestParseEncryptedToken(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		assertProvider string
		assertValue    string
		assertError    bool
	}

	run := func(t *testing.T, tc testCase) {
		got, err := ParseEncryptedToken(crypto.Data(tc.input), 1, 2)

		if tc.assertError {
			assert.Assert(t, err != nil)
		} else {
			assert.NilError(t, err)
			assert.Assert(t, got != nil)

			assert.Equal(t, got.Provider(), tc.assertProvider)
			assert.Equal(t, got.Provider(), tc.assertProvider)
			assert.Equal(t, got.Value().String(), tc.assertValue)
		}
	}
	testCases := []testCase{
		{
			name:           "empty plain",
			input:          "ENC[TEST]",
			assertProvider: "TEST",
			assertValue:    "",
		},
		{
			name:           "common usage",
			input:          "ENC[TEST,data]",
			assertProvider: "TEST",
			assertValue:    "data",
		},
		{
			name:           "trailing delim",
			input:          "ENC[TEST,]",
			assertProvider: "TEST",
			assertValue:    "",
		},
		{
			name:        "empty provider",
			input:       "ENC[,data]",
			assertValue: "data",
		},
		{
			name:           "trim params",
			input:          "ENC[ TEST , fir\tst \r.\n sec\vond ]",
			assertProvider: "TEST",
			assertValue:    "first.second",
		},
		{
			name:        "too many params",
			input:       "ENC[TEST,data,error]",
			assertError: true,
		},
		{
			name:        "malformed opening",
			input:       "ENC<TEST,data>",
			assertError: true,
		},
		{
			name:        "too short",
			input:       "ENC[",
			assertError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func TestParseEncryptedValue(t *testing.T) {
	type testCase struct {
		name      string
		input     string
		assertion string
	}

	run := func(t *testing.T, tc testCase) {
		got := ParseEncryptedValue([]byte(tc.input))

		assert.Equal(t, string(got), tc.assertion)
	}
	testCases := []testCase{
		{
			name: "empty",
		},
		{
			name:  "only spaces",
			input: "\t \r \v \n",
		},
		{
			name:      "leading space",
			input:     "\ttest",
			assertion: "test",
		},
		{
			name:      "trailing space",
			input:     "test\r",
			assertion: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
