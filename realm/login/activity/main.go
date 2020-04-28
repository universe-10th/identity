package activity

import (
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/credentials/traits/deniable"
	"github.com/universe-10th/identity/realm"
)

// This pipeline step tells when a credential could not
// login because it counts as inactive.
type ActivityStep uint8

// Attempts a log-in step which would fail if the credential
// counts as inactive.
func (ActivityStep) Login(credential credentials.Credential, password string) error {
	if activable, ok := credential.(deniable.Activable); ok && !activable.Active() {
		return realm.ErrLoginFailed
	} else {
		return nil
	}
}
