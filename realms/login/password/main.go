package password

import (
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/realms"
)

// This step checks the password for the given
// credential, using its hasher. It may return
// errors of invalid password or of the credential
// not being able to login because it has none.
type PasswordCheckingStep uint8

// Attempts the login step of password check.
func (PasswordCheckingStep) Login(credential credentials.Credential, password string) error {
	hashed := credential.HashedPassword()
	if hashed == "" {
		return realms.ErrLoginFailed
	}

	hasher := credential.Hasher()
	if err := hasher.Validate(password, credential.HashedPassword()); err != nil {
		return realms.ErrLoginFailed
	} else {
		return nil
	}
}
