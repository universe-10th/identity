package compound

import (
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/authreqs/scoped"
	"github.com/universe-10th/identity/authreqs/superuser"
	"github.com/universe-10th/identity/authreqs/staff"
	"github.com/universe-10th/identity/authreqs"
)

// Returns an authorizer function which tells
// whether the credential was authorized by
// trying different authorization requirements
// until one succeeds.
func TryAll(alternatives ...authreqs.AuthorizationRequirement) func(credentials.Credential) bool {
	return func(credential credentials.Credential) bool {
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
func Admin(scopes ...scoped.ScopeSpec) func(credentials.Credential) bool {
	superuserReq := superuser.RequireSuperuser
	staffReq := staff.RequireStaff
	scopedReq := scoped.RequireScopesAmong(scopes...)

	return func(credential credentials.Credential) bool {
		if superuserReq.SatisfiedBy(credential) {
			return true
		} else {
			return staffReq.SatisfiedBy(credential) && scopedReq.SatisfiedBy(credential)
		}
	}
}