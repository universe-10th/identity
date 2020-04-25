package authreqs

// An authorization requirement is a criteria being
// tested against a particular credential.
type AuthorizationRequirement interface {
	// A mean to test whether a particular credential
	// satisfies this requirement.
	SatisfiedBy(identity.Credential) bool
}
