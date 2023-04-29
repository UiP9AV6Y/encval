package crypto

import (
	"strings"
)

// Encrypters is a map implementation with convenience functions
type Encrypters map[string]Encrypter

// Add stores the given encrypter implementation under the given key.
// Keys can only be used once and return false on collisions.
// The storage key is sanitized for unambiguous access later on.
func (e Encrypters) Add(provider string, impl Encrypter) bool {
	needle := strings.ToUpper(provider)
	if _, dup := e[needle]; dup {
		return false
	}

	e[needle] = impl

	return true
}

// Get retrieves a previously stored encrypter implementation
// using the given key as lookup. If no instance is available,
// false is returned.
func (e Encrypters) Get(provider string) (Encrypter, bool) {
	needle := strings.ToUpper(provider)
	result, ok := e[needle]

	return result, ok
}

// Providers returns all the keys of the map.
func (e Encrypters) Providers() []string {
	keys := make([]string, len(e))
	i := 0

	for k := range e {
		keys[i] = k
		i++
	}

	return keys
}

// Encrypter is a processing contract for data in- and output.
// It is up to the implementation to decide if they can consume
// their own output multiple times.
type Encrypter interface {
	Encrypt(data Data) (Data, error)
	Decrypt(data Data) (Data, error)
}
