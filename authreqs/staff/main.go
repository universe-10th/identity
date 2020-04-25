package staff

import (
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/credentials/traits/staff"
)

// This requirement checks whether a
// credential is staff user. It is meant
// to be used in a "TryAll" compound auth
// requirement and in the first position
// of such array. It will be used by a
// convenience class for the developer.
type Staff uint8

// Tests whether the credential is staff or not.
func (Staff) SatisfiedBy(credential credentials.Credential) bool {
	if capable, ok := credential.(staff.StaffCapable); ok && capable.Staff() {
		return true
	} else {
		return false
	}
}

const RequireStaff = Staff(0)