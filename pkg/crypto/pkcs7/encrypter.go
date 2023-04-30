package pkcs7

import (
	"context"
	"crypto/rsa"
	"crypto/x509"

	"go.mozilla.org/pkcs7"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
	"github.com/UiP9AV6Y/encval/pkg/log"
)

const Encrypter = "PKCS7"

type PKCS7Encrypter struct {
	key  *PrivateKey
	cert *PublicKey

	keyCache  *rsa.PrivateKey
	certCache *x509.Certificate
}

func New(key *PrivateKey, cert *PublicKey) *PKCS7Encrypter {
	result := &PKCS7Encrypter{
		key:  key,
		cert: cert,
	}

	return result
}

func (e *PKCS7Encrypter) Encrypt(data crypto.Data) (crypto.Data, error) {
	if e.certCache == nil {
		cert, err := e.cert.Load()
		if err != nil {
			return nil, err
		}
		e.certCache = cert
	}

	// use the same algorithm as hiera-eyaml
	pkcs7.ContentEncryptionAlgorithm = pkcs7.EncryptionAlgorithmAES256CBC
	payload, err := pkcs7.Encrypt(data, []*x509.Certificate{e.certCache})
	if err != nil {
		return nil, err
	}

	return crypto.Encode(payload)
}

func (e *PKCS7Encrypter) Decrypt(data crypto.Data) (crypto.Data, error) {
	if e.keyCache == nil {
		key, err := e.key.Load()
		if err != nil {
			return nil, err
		}
		e.keyCache = key
	}

	if e.certCache == nil {
		cert, err := e.cert.Load()
		if err != nil {
			return nil, err
		}
		e.certCache = cert
	}

	raw, err := crypto.Decode(data)
	if err != nil {
		return nil, err
	}

	p7, err := pkcs7.Parse(raw)
	if err != nil {
		return nil, err
	}

	payload, err := p7.Decrypt(e.certCache, e.keyCache)
	if err != nil {
		return nil, err
	}

	return crypto.Data(payload), nil
}

func (p *PKCS7Encrypter) GenerateSecrets(force bool, ctx context.Context) error {
	logger, err := log.FromContext(ctx)
	if err != nil {
		return err
	}

	key, err := p.key.Generate()
	if err != nil {
		return err
	}

	cert, err := p.cert.Generate(key)
	if err != nil {
		return err
	}

	if err := p.cert.Save(cert, force); err != nil {
		return err
	}

	logger.Info().Printf("Saving %s cert for %q using %v in %s\n",
		Encrypter,
		p.cert.Subject(),
		p.cert.SignatureAlgorithm(),
		p.cert.Path())

	if err := p.key.Save(key, force); err != nil {
		return err
	}

	logger.Info().Printf("Saving %s key of size %d in %s\n",
		Encrypter,
		p.key.Bits(),
		p.key.Path())

	return nil
}
