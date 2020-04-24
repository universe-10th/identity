package old

// An authorization requirement is a criteria being
// tested against a particular credential.
type AuthorizationRequirement interface {
	// A mean to test whether a particular credential
	// satisfies this requirement.
	SatisfiedBy(Credential) bool
}

// A scope is a particular requirement having a unique
// identification. This means: a scope is an atomic
// authorization requirement. A scope is often
// referred as a Permission.
type Scope interface {
	AuthorizationRequirement
	// The fully-qualified scope key.
	Key() string
	// The scope friendly name.
	Name() string
	// The scope description (optional).
	Description() string
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
