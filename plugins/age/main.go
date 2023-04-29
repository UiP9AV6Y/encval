package main

import (
	"path/filepath"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

type agePlugin struct {
	Options struct {
		IdentitiesDir string `kong:"help='Directory containing AGE identities',type='path',placeholder='DIR'"`
		CertDir       string `kong:"help='Directory containing AGE certificates',type='path',placeholder='DIR'"`
	} `kong:"embed,prefix='age.'"`
}

func (p *agePlugin) Encrypter() string {
	return Encrypter
}

func (p *agePlugin) NewEncrypter(baseDir string) (crypto.Encrypter, error) {
	priv := p.Options.IdentitiesDir
	if priv == "" {
		priv = filepath.Join(baseDir, "age")
	}

	pub := p.Options.CertDir
	if pub == "" {
		pub = filepath.Join(baseDir, "age")
	}

	key := Identities(priv)
	cert := Recipients(pub)

	return New(key, cert), nil
}

func NewEncrypterPlugin() crypto.EncrypterPlugin {
	result := &agePlugin{}

	return result
}
