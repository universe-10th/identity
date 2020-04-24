package deniable

import (
	"time"
	"github.com/universe-10th/identity"
)

// This trait allows to punish a credential (this implies the
// credential should fail to log in for certain duration or
// perhaps permanently). If the first return value is nil,
// the credential is not punished. If not, the third return
// value will hold the reason. Also, if not nil and the second
// value is nil, this means a permanent punish. The fourth result
// value may be nil. Punishing an already punished credential
// will entirely replace the ban.
type Punishable interface {
	PunishedFor() (punishedOn *time.Time, forTime *time.Duration, reason interface{}, by identity.Credential)
	Punish(forTime *time.Duration, reason interface{}, by identity.Credential)
	Unpunish()
}

// This trait allows to inactivate a credential. Inactive
// credentials will act as if they do not exist regarding a
// login operation.
type Activable interface {
	Active() bool
}