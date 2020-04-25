package authreqs

import "github.com/universe-10th/identity/credentials"

// An authorization requirement is a criteria being
// tested against a particular credential.
type AuthorizationRequirement interface {
	// A mean to test whether a particular credential
	// satisfies this requirement.
	SatisfiedBy(credentials.Credential) bool
}
