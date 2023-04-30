package pkcs7

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPrivateKey(t *testing.T) {
	f, err := os.CreateTemp("", "example")
	assert.NilError(t, err)
	defer os.Remove(f.Name())
	assert.NilError(t, f.Close())

	out := NewPrivateKey(f.Name(), 2048)

	key, err := out.Generate()
	assert.NilError(t, err)

	err = out.Save(key, false) // tmpfile already exists
	assert.Assert(t, err != nil)

	err = out.Save(key, true)
	assert.NilError(t, err)

	key2, err2 := out.Load()
	assert.NilError(t, err2)
	assert.Assert(t, key.Equal(key2))
}
