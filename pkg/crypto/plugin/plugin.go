package plugin

import (
	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

// Plugin is a crypto.EncrypterPlugin implementation using an dynamically
// loaded implementation.
type Plugin struct {
	path string

	crypto.EncrypterPlugin `kong:"embed,group='encryption',prefix='plugin.'"`
}

// String returns the filesystem location the plugin was loaded from
func (p Plugin) String() string {
	return p.path
}
