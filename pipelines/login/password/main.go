package password

import (
	"github.com/universe-10th/identity"
	"github.com/kataras/iris/core/errors"
)

// This step checks the password for the given
// credential, using its hasher. It may return
// errors of invalid password or of the credential
// not being able to login because it has none.
type PasswordCheckingStep uint8


var ErrInvalidPassword = errors.New("invalid password")
var ErrNoPassword = errors.New("unset password to check against")

func (PasswordCheckingStep) Login(credential identity.Credential, password string) error {
	hashed := credential.HashedPassword()
	if hashed == "" {
		return ErrNoPassword
	}

	hasher := credential.Engine()
	if err := hasher.Validate(password, credential.HashedPassword()); err != nil {
		return ErrInvalidPassword
	} else {
		return nil
	}
}

