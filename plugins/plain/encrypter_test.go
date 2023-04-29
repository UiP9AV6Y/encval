package main

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

func TestEncrypt(t *testing.T) {
	type testCase struct {
		name         string
		input        string
		urlEncoding  bool
		padding      bool
		assertError  bool
		assertOutput string
	}

	run := func(t *testing.T, tc testCase) {
		out = New(tc.urlEncoding, tc.padding)
		assert.Assert(t, out != nil)

		got, err := out.Encrypt(crypto.Data(tc.input))

		if tc.assertError {
			assert.Assert(t, err != nil)
		} else {
			assert.Equal(t, string(got), tc.assertOutput)
		}
	}
	testCases := []testCase{
		{
			name:         "empty",
			input:        "",
			assertOutput: "",
		},
		{
			name:         "empty (URL encoding)",
			input:        "",
			assertOutput: "",
			urlEncoding:  true,
		},
		{
			name:         "space",
			input:        " ",
			assertOutput: "IA",
		},
		{
			name:         "test",
			input:        `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua`,
			assertOutput: `TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdCwgc2VkIGRvIGVpdXNtb2QgdGVtcG9yIGluY2lkaWR1bnQgdXQgbGFib3JlIGV0IGRvbG9yZSBtYWduYSBhbGlxdWE`,
		},
		{
			name:         "test (padding)",
			input:        `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua`,
			assertOutput: `TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdCwgc2VkIGRvIGVpdXNtb2QgdGVtcG9yIGluY2lkaWR1bnQgdXQgbGFib3JlIGV0IGRvbG9yZSBtYWduYSBhbGlxdWE=`,
			padding:      true,
		},
		{
			name:         "test (padding, URL encoding)",
			input:        `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua`,
			assertOutput: `TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdCwgc2VkIGRvIGVpdXNtb2QgdGVtcG9yIGluY2lkaWR1bnQgdXQgbGFib3JlIGV0IGRvbG9yZSBtYWduYSBhbGlxdWE=`,
			urlEncoding:  true,
			padding:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
