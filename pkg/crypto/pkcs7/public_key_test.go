package pkcs7

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPublicKey(t *testing.T) {
	f, err := os.CreateTemp("", "example")
	assert.NilError(t, err)
	defer os.Remove(f.Name())
	assert.NilError(t, f.Close())

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NilError(t, err)

	out := NewPublicKey(f.Name(), "spec.test", x509.SHA256WithRSA)

	cert, err := out.Generate(key)
	assert.NilError(t, err)
	assert.Equal(t, cert.SignatureAlgorithm, x509.SHA256WithRSA)

	err = out.Save(cert, false) // tmpfile already exists
	assert.Assert(t, err != nil)

	err = out.Save(cert, true)
	assert.NilError(t, err)

	cert2, err2 := out.Load()
	assert.NilError(t, err2)
	assert.Equal(t, cert2.SerialNumber.String(), cert.SerialNumber.String())
}
