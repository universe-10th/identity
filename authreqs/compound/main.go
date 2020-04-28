package compound

import (
	"github.com/universe-10th/identity/authreqs"
	"github.com/universe-10th/identity/authreqs/scoped"
	"github.com/universe-10th/identity/authreqs/staff"
	"github.com/universe-10th/identity/authreqs/superuser"
	"github.com/universe-10th/identity/credentials"
)

// A requirement that tests a list of requirements
// and succeeds when one of them succeeds.
type TryAll []authreqs.AuthorizationRequirement

// Performs the test one-by-one until one succeeds.
func (tryAll TryAll) SatisfiedBy(credential credentials.Credential) bool {
	for _, alternative := range tryAll {
		if alternative.SatisfiedBy(credential) {
			return true
		}
	}
	return false
}

// Intended for admin side, it tests whether the
// credential is superuser, or a staff satisfying
// the required permissions, if any.
type AdminRequirement struct {
	superuserReq superuser.Superuser
	staffReq     staff.Staff
	scopedReq    *scoped.Scoped
}

// Performs the test (superuser || staff && permissions).
func (adminRequirement *AdminRequirement) SatisfiedBy(credential credentials.Credential) bool {
	if adminRequirement.superuserReq.SatisfiedBy(credential) {
		return true
	} else {
		return adminRequirement.staffReq.SatisfiedBy(credential) && adminRequirement.scopedReq.SatisfiedBy(credential)
	}
}

// Returns a superuser/scoped combined authorizer
// which tells whether the credential is superuser.
// This is not useful for user-end features but
// only for administrative features.
func Admin(scopes ...scoped.ScopeSpec) *AdminRequirement {
	superuserReq := superuser.RequireSuperuser
	staffReq := staff.RequireStaff
	scopedReq := scoped.RequireScopesAmong(scopes...)
	return &AdminRequirement{superuserReq, staffReq, scopedReq}
}
