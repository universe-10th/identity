package superuser

import (
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/credentials/traits/superuser"
)

// This requirement checks whether a
// credential is superuser. It is meant
// to be used in an "TryAll" compound auth
// requirement and in the first position
// of such array. It will be used by a
// convenience class for the developer.
type Superuser uint8

// Tests whether the credential is superuser or not.
func (Superuser) SatisfiedBy(credential credentials.Credential) bool {
	if capable, ok := credential.(superuser.SuperuserCapable); ok && capable.Superuser() {
		return true
	} else {
		return false
	}
}

const RequireSuperuser = Superuser(0)