package pkcs7

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/UiP9AV6Y/encval/pkg/fs"
)

type PublicKey struct {
	path    string
	subject string
	algo    x509.SignatureAlgorithm
}

func NewPublicKey(path, subject string, algo x509.SignatureAlgorithm) *PublicKey {
	result := &PublicKey{
		path:    path,
		subject: subject,
		algo:    algo,
	}

	return result
}

func (k *PublicKey) SignatureAlgorithm() x509.SignatureAlgorithm {
	return k.algo
}

func (k *PublicKey) Subject() string {
	return k.subject
}

func (k *PublicKey) Path() string {
	return k.path
}

func (k *PublicKey) Generate(key *rsa.PrivateKey) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 32)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	sub := pkix.Name{
		CommonName: k.subject,
	}
	csr := &x509.Certificate{
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(50, 0, 0),
		SerialNumber:          serialNumber,
		Version:               2,
		Subject:               sub,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		SignatureAlgorithm:    k.algo,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, csr, csr, key.Public(), key)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(cert)
}

func (k *PublicKey) Save(cert *x509.Certificate, force bool) error {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	data := pem.EncodeToMemory(block)

	if force {
		return fs.ReinstallFile(k.path, data, 0644)
	}

	return fs.InstallFile(k.path, data, 0644)
}

func (k *PublicKey) Load() (*x509.Certificate, error) {
	data, err := os.ReadFile(k.path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("Failed to parse public key %q", k.path)
	}

	return x509.ParseCertificate(block.Bytes)
}
