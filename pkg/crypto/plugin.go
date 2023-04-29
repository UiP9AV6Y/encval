package crypto

import ()

// EncrypterPlugin is a factory contract for encryption systems.
type EncrypterPlugin interface {
	// Encrypter is an identification method. The return value is
	// expected to be stable over multiple invocations.
	Encrypter() string
	// NewEncrypter is a factory method. The input parameter
	// is a directory for potential data persistence, which might not exist.
	NewEncrypter(string) (Encrypter, error)
}
