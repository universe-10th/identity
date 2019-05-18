package utils

import (
	"errors"
	"github.com/luismasuelli/go-identity/stub"
)

var Unauthorized = errors.New("unauthorized to execute the action")

/**
 * Checks whether the credential is authorized. The given credential must be
 *   a (*T) value (with T being a struct type).
 */
func Authorize(credential stub.Credential, requirement stub.AuthorizationRequirement) error {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(credential) {
		return StructPointerStubExpected
	}

	// Empty authorization is always true
	if requirement == nil {
		return nil
	}

	if requirement.SatisfiedBy(credential) {
		return nil
	} else {
		return Unauthorized
	}
}
