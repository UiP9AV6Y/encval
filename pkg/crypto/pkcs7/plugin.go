package pkcs7

import (
	"crypto/x509"
	"path/filepath"
)

type PKCS7Plugin struct {
	PrivateKey string `kong:"help='Path to private key',type='path',placeholder='FILE'"`
	PublicKey  string `kong:"help='Path to public key',type='path',placeholder='FILE'"`
	Subject    string `kong:"help='Subject to use for certificate when creating keys',default='/'"`
	Digest     string `kong:"help='Hash function used for PKCS7',enum='SHA256,SHA384,SHA512',default='SHA256'"`
	Keysize    uint32 `kong:"help='Key size used for encryption',default='2048'"`
}

func (o *PKCS7Plugin) Encrypter() string {
	return Encrypter
}

func (o *PKCS7Plugin) NewEncrypter(baseDir string) *PKCS7Encrypter {
	var algo x509.SignatureAlgorithm
	switch o.Digest {
	case "SHA384":
		algo = x509.SHA384WithRSA
	case "SHA512":
		algo = x509.SHA512WithRSA
	default:
		algo = x509.SHA256WithRSA
	}

	priv := o.PrivateKey
	if priv == "" {
		priv = filepath.Join(baseDir, "pkcs7", "private_key.pkcs7.pem")
	}

	pub := o.PublicKey
	if pub == "" {
		pub = filepath.Join(baseDir, "pkcs7", "public_key.pkcs7.pem")
	}

	key := NewPrivateKey(priv, o.Keysize)
	cert := NewPublicKey(pub, o.Subject, algo)

	return New(key, cert)
}
