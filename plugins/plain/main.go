package main

import (
	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

type plainPlugin struct {
	Options struct {
		UrlEncode bool `kong:"help='Use URL-safe encoding',negatable"`
		Padding   bool `kong:"help='Include padding characters',negatable"`
	} `kong:"embed,prefix='plain.'"`
}

func (p *plainPlugin) Encrypter() string {
	return Encrypter
}

func (p *plainPlugin) NewEncrypter(_ string) (crypto.Encrypter, error) {
	return New(p.Options.UrlEncode, p.Options.Padding), nil
}

func NewEncrypterPlugin() crypto.EncrypterPlugin {
	result := &plainPlugin{}

	return result
}
