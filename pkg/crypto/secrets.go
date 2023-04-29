package crypto

import (
	"context"
)

// SecretsGenerator is implemented by encryption plugins which utilize
// some kind of persistent information during processing.
type SecretsGenerator interface {
	// GenerateSecrets is used to create any configuration, secrets, or otherwise
	// persistent data for later de-/encryption. If the first parameter is true,
	// any existing data is expected to be overwritten. The given context contains
	// additional utilities.
	GenerateSecrets(bool, context.Context) error
}
