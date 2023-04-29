package parser

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

func TestParseDecryptedToken(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		assertIndex    int
		assertProvider string
		assertValue    string
		assertError    bool
	}

	run := func(t *testing.T, tc testCase) {
		got, err := ParseDecryptedToken(crypto.Data(tc.input), 1, 2)

		if tc.assertError {
			assert.Assert(t, err != nil)
		} else {
			assert.NilError(t, err)
			assert.Assert(t, got != nil)

			assert.Equal(t, got.Index, tc.assertIndex)
			assert.Equal(t, got.Provider(), tc.assertProvider)
			assert.Equal(t, got.Value().String(), tc.assertValue)
		}
	}
	testCases := []testCase{
		{
			name:           "empty plain",
			input:          "DEC::TEST[]!",
			assertProvider: "TEST",
		},
		{
			name:           "existing empty",
			input:          "DEC(1)::TEST[]!",
			assertProvider: "TEST",
			assertIndex:    1,
		},
		{
			name:        "empty existing provider",
			input:       "DEC(1)::[]!",
			assertIndex: 1,
		},
		{
			name:  "empty new provider",
			input: "DEC::[]!",
		},
		{
			name:           "existing value",
			input:          "DEC(1)::TEST[test]!",
			assertProvider: "TEST",
			assertValue:    "test",
			assertIndex:    1,
		},
		{
			name:           "new value",
			input:          "DEC::TEST[test]!",
			assertProvider: "TEST",
			assertValue:    "test",
		},
		{
			name:           "lenient parsing",
			input:          "XXC::X[]!",
			assertProvider: "X",
		},
		{
			name:        "malformed index",
			input:       "DEC)1(::TEST[test]!",
			assertError: true,
		},
		{
			name:        "missing index",
			input:       "DEC()::TEST[test]!",
			assertError: true,
		},
		{
			name:        "index without rest",
			input:       "DEC(1234567890[])",
			assertError: true,
		},
		{
			name:        "negative index",
			input:       "DEC(-1)::TEST[test]!",
			assertError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
