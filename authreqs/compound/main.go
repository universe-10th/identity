package compound

import "github.com/universe-10th/identity"

// This requirement implements satisfaction
// by checking and requiring all of the composed
// (children) requirements.
type AllOf []identity.AuthorizationRequirement

// Considers satisfaction by requiring all the
// inner requirements to be satisfied.
func (allOf AllOf) SatisfiedBy(credential identity.Credential) bool {
	for _, value := range allOf {
		if !value.SatisfiedBy(credential) {
			return false
		}
	}
	return true
}

// This requirement implements satisfaction
// by checking and requiring all of the composed
// (children) requirements.
type AnyOf []identity.AuthorizationRequirement

// Considers satisfaction by requiring any of
// the inner requirements to be satisfied.
func (anyOf AnyOf) SatisfiedBy(credential identity.Credential) bool {
	for _, value := range anyOf {
		if value.SatisfiedBy(credential) {
			return true
		}
	}
	return false
}
