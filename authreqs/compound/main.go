package compound

import (
	"github.com/universe-10th/identity"
	"github.com/universe-10th/identity/authreqs/scoped"
	"github.com/universe-10th/identity/authreqs/superuser"
	"github.com/universe-10th/identity/authreqs/staff"
)

// Returns an authorizer function which tells
// whether the credential was authorized by
// trying different authorization requirements
// until one succeeds.
func TryAll(alternatives ...identity.AuthorizationRequirement) func(identity.Credential) bool {
	return func(credential identity.Credential) bool {
		for _, alternative := range alternatives {
			if alternative.SatisfiedBy(credential) {
				return true
			}
		}
		return false
	}
}

// Returns a superuser/scoped combined authorizer
// which tells whether the credential is superuser.
// This is not useful for user-end features but
// only for administrative features.
func Admin(scopes ...scoped.ScopeSpec) func(identity.Credential) bool {
	superuserReq := superuser.RequireSuperuser
	staffReq := staff.RequireStaff
	scopedReq := scoped.RequireScopesAmong(scopes...)

	return func(credential identity.Credential) bool {
		if superuserReq.SatisfiedBy(credential) {
			return true
		} else {
			return staffReq.SatisfiedBy(credential) && scopedReq.SatisfiedBy(credential)
		}
	}
}