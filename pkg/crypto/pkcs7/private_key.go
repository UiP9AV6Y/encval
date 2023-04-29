package pkcs7

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/UiP9AV6Y/encval/pkg/fs"
)

type PrivateKey struct {
	path string
	bits uint32
}

func NewPrivateKey(path string, bits uint32) *PrivateKey {
	result := &PrivateKey{
		path: path,
		bits: bits,
	}

	return result
}

func (k *PrivateKey) Bits() uint32 {
	return k.bits
}

func (k *PrivateKey) Path() string {
	return k.path
}

func (k *PrivateKey) Generate() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, int(k.bits))
}

func (k *PrivateKey) Save(key *rsa.PrivateKey, force bool) error {
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	data := pem.EncodeToMemory(block)

	if force {
		return fs.ReinstallFile(k.path, data, 0600)
	}

	return fs.InstallFile(k.path, data, 0600)
}

func (k *PrivateKey) Load() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(k.path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("Failed to parse private key %q", k.path)
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
