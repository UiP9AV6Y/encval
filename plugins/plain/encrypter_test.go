package main

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

func TestEncrypt(t *testing.T) {
	type testCase struct {
		name        string
		input       string
		urlEncoding bool
		padding     bool
	}

	run := func(t *testing.T, tc testCase) {
		out := New(tc.urlEncoding, tc.padding)
		assert.Assert(t, out != nil)

		enc, err := out.Encrypt(crypto.Data(tc.input))
		assert.NilError(t, err)
		dec, err := out.Decrypt(enc)
		assert.NilError(t, err)

		assert.Equal(t, string(dec), tc.input)
	}
	testCases := []testCase{
		{
			name:  "empty",
			input: "",
		},
		{
			name:        "empty (URL encoding)",
			input:       "",
			urlEncoding: true,
		},
		{
			name:  "space",
			input: " ",
		},
		{
			name:  "test",
			input: `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua`,
		},
		{
			name:    "test (padding)",
			input:   `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua`,
			padding: true,
		},
		{
			name:        "test (padding, URL encoding)",
			input:       `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua`,
			urlEncoding: true,
			padding:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
