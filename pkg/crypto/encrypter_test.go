package crypto

import (
	"testing"

	"gotest.tools/v3/assert"
)

type SpecEncrypter struct {
	ID int
}

func (e *SpecEncrypter) Encrypt(data Data) (Data, error) {
	return Data([]byte{}), nil
}

func (e *SpecEncrypter) Decrypt(data Data) (Data, error) {
	return Data([]byte{}), nil
}

func TestEncrypters(t *testing.T) {
	var conv *SpecEncrypter
	var got Encrypter
	var ok bool

	out := make(Encrypters, 2)
	enc := &SpecEncrypter{ID: 1}

	got, ok = out.Get("spec")
	assert.Assert(t, !ok)

	ok = out.Add("spec", enc)
	assert.Assert(t, ok)

	got, ok = out.Get("spec")
	assert.Assert(t, ok)

	conv, ok = got.(*SpecEncrypter)
	assert.Assert(t, ok)
	assert.Assert(t, conv.ID == 1)

	got, ok = out.Get("SPEC")
	assert.Assert(t, ok)

	ok = out.Add("SpEc", enc)
	assert.Assert(t, !ok)

	ok = out.Add("SpEc", &SpecEncrypter{ID: 2})
	assert.Assert(t, !ok)

	got, ok = out.Get("SpEc")
	assert.Assert(t, ok)

	conv, ok = got.(*SpecEncrypter)
	assert.Assert(t, ok)
	assert.Assert(t, conv.ID == 1)

	prov := out.Providers()
	assert.Equal(t, 1, len(prov))
	assert.Equal(t, prov[0], "SPEC")
}
