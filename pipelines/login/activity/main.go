package activity

import (
	"github.com/universe-10th/identity"
	"github.com/universe-10th/identity/traits/credential/deniable"
	"github.com/universe-10th/identity/pipelines/login"
)

// This pipeline step tells when a credential could not
// login because it counts as inactive.
type ActivityStep uint8

// Attempts a log-in step which would fail if the credential
// counts as inactive.
func (ActivityStep) Login(credential identity.Credential, password string) error {
	if activable, ok := credential.(deniable.Activable); ok && !activable.Active() {
		return login.ErrLoginFailed
	} else {
		return nil
	}
}
