package types

import "github.com/luismasuelli/go-identity/interfaces"

type AllOf []interfaces.AuthorizationRequirement
type AnyOf []interfaces.AuthorizationRequirement


func (allOf AllOf) SatisfiedBy(credential interfaces.Credential) bool {
	for _, value := range allOf {
		if !value.SatisfiedBy(credential) {
			return false
		}
	}
	return true
}


func (anyOf AnyOf) SatisfiedBy(credential interfaces.Credential) bool {
	for _, value := range anyOf {
		if value.SatisfiedBy(credential) {
			return true
		}
	}
	return false
}