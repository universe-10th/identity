package identity

// An authorization requirement is a criteria being
// tested against a particular credential.
type AuthorizationRequirement interface {
	// A mean to test whether a particular credential
	// satisfies this requirement.
	SatisfiedBy(Credential) bool
}

type AllOf []AuthorizationRequirement
type AnyOf []AuthorizationRequirement

func (allOf AllOf) SatisfiedBy(credential Credential) bool {
	for _, value := range allOf {
		if !value.SatisfiedBy(credential) {
			return false
		}
	}
	return true
}

func (anyOf AnyOf) SatisfiedBy(credential Credential) bool {
	for _, value := range anyOf {
		if value.SatisfiedBy(credential) {
			return true
		}
	}
	return false
}
