package plugin

import (
	"github.com/UiP9AV6Y/encval/pkg/crypto"
	"github.com/UiP9AV6Y/encval/pkg/crypto/pkcs7"
)

// PKCS7EncryptionPlugin is a CLI argument container
type PKCS7EncryptionPlugin struct {
	*pkcs7.PKCS7Plugin `kong:"embed,group='encryption',prefix='pkcs7.'"`
}

func (b *PKCS7EncryptionPlugin) Encrypter() string {
	return b.PKCS7Plugin.Encrypter()
}

func (b *PKCS7EncryptionPlugin) NewEncrypter(configDir string) (crypto.Encrypter, error) {
	return b.PKCS7Plugin.NewEncrypter(configDir), nil
}

func NewDefaultPlugin() crypto.EncrypterPlugin {
	return &PKCS7EncryptionPlugin{}
}
