package activity

import (
	"github.com/universe-10th/identity"
	"github.com/universe-10th/identity/traits/credential/deniable"
	"errors"
)

var ErrInactive = errors.New("credential is inactive")

// This pipeline step tells when a credential could not
// login because it counts as inactive.
type ActivityStep uint8

// Attempts a log-in step which would fail if the credential
// counts as inactive.
func (ActivityStep) Login(credential identity.Credential, password string) error {
	if activable, ok := credential.(deniable.Activable); ok && !activable.Active() {
		return ErrInactive
	} else {
		return nil
	}
}
