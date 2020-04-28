package credentials

import (
	"github.com/universe-10th/identity/hashing"
)

// Credentials have a password that has to be
// matched. They also provide a hashing engine
// to be used, both to set the password and to
// validate a password against the hash.
type Credential interface {
	// Returns the current (hashed) password.
	HashedPassword() string
	// Sets a new (already hashed) password.
	SetHashedPassword(string)
	// Class method to return the hashing
	// engine used to hash and validate the
	// passwords.
	Hasher() hashing.HashingEngine
}
