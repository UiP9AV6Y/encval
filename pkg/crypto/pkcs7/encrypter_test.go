package pkcs7

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

func TestEncrypt(t *testing.T) {
	type testCase struct {
		name        string
		input       string
		assertError bool
	}

	priv := &PrivateKey{
		path: "testdata/insecure.key.pem",
	}
	pub := &PublicKey{
		path: "testdata/insecure.crt.pem",
	}
	run := func(t *testing.T, tc testCase) {
		out := New(priv, pub)
		assert.Assert(t, out != nil)

		got, err := out.Encrypt(crypto.Data(tc.input))

		if tc.assertError {
			assert.Assert(t, err != nil)
		} else {
			// unless we fake the encryption, the result is not stable
			assert.Assert(t, len(got) > 0)
		}
	}
	testCases := []testCase{
		{
			name:  "empty",
			input: "",
		},
		{
			name:  "space",
			input: " ",
		},
		{
			name:  "multiline",
			input: "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut labore et dolore magna aliqua",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
