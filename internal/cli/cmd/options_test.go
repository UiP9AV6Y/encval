package cmd

import (
	"testing"

	"github.com/alecthomas/kong"
	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

type testEncrypter string

func (e testEncrypter) Encrypt(data crypto.Data) (crypto.Data, error) {
	return crypto.Data(e), nil
}

func (e testEncrypter) Decrypt(data crypto.Data) (crypto.Data, error) {
	return crypto.Data(e), nil
}

type testPlugin string

func (p testPlugin) Encrypter() string {
	return string(p)
}

func (p testPlugin) NewEncrypter(dir string) (crypto.Encrypter, error) {
	return testEncrypter(dir), nil
}

func TestOptionsNewEncrypters(t *testing.T) {
	out := NewGlobalOptions()
	out.Config = kong.ConfigFlag("testdata/config.toml")
	out.Plugins = kong.Plugins{
		testPlugin("TEST"),
	}

	got, err := out.NewEncrypters()
	assert.NilError(t, err)
	assert.Equal(t, len(got), 1)

	enc, ok := got.Get("test")
	assert.Assert(t, enc != nil)
	assert.Equal(t, ok, true)

	dir, err := enc.Decrypt(crypto.Data(""))
	assert.NilError(t, err)
	assert.Equal(t, string(dir), "testdata")
}

func TestStringSliceContains(t *testing.T) {
	type testCase struct {
		name         string
		input        []string
		needle       string
		assertOutput bool
	}

	run := func(t *testing.T, tc testCase) {
		got := stringSliceContains(tc.input, tc.needle)

		assert.Equal(t, got, tc.assertOutput)
	}
	testCases := []testCase{
		{
			name:  "empty",
			input: []string{},
		},
		{
			name:         "one",
			input:        []string{""},
			assertOutput: true,
		},
		{
			name:         "two",
			input:        []string{"", ""},
			assertOutput: true,
		},
		{
			name:   "missed",
			input:  []string{"", ""},
			needle: "test",
		},
		{
			name:         "last",
			input:        []string{"1", "2", "3", "4"},
			needle:       "4",
			assertOutput: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
