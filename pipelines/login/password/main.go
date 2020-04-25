package password

import (
	"github.com/universe-10th/identity"
	"github.com/universe-10th/identity/pipelines/login"
)

// This step checks the password for the given
// credential, using its hasher. It may return
// errors of invalid password or of the credential
// not being able to login because it has none.
type PasswordCheckingStep uint8


func (PasswordCheckingStep) Login(credential identity.Credential, password string) error {
	hashed := credential.HashedPassword()
	if hashed == "" {
		return login.ErrLoginFailed
	}

	hasher := credential.Engine()
	if err := hasher.Validate(password, credential.HashedPassword()); err != nil {
		return login.ErrLoginFailed
	} else {
		return nil
	}
}

