package identity

import (
	"errors"
	"reflect"
	"github.com/universe-10th/identity/support/types"
	"github.com/universe-10th/identity/stub"
)

var StructPointerStubExpected = errors.New("only pointer-kind stubs are allowed")
var CredentialDoesNotHavePassword = errors.New("fetched credential does not have password")

func prototypeIsAStructPtr(prototype interface{}) bool {
	// nil interface is not allowed
	if prototype == nil {
		return false
	}
	// Only prototypes that are a pointer types are allowed.
	prType := reflect.TypeOf(prototype)
	if prType.Kind() != reflect.Ptr {
		return false
	}
	// Also the indirect type must be a struct.
	return prType.Elem().Kind() != reflect.Struct
}


// Given a database query, a model prototype, an identification and a password, it tries to
// perform a login. It may fail for several reasons: the prototype is not a (*T) type (with
// T being a struct type), a database error, a bad password, a password-less credential, or
// another custom login error (after or before the passwords check).
func Login(source stub.Source, lookupResult stub.Credential, identification interface{}, password string) error {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(lookupResult) {
		return StructPointerStubExpected
	}

	// Find the credential in database
	if err := source.ByIdentification(lookupResult, identification); err != nil {
		return err
	}

	hash := lookupResult.HashedPassword()

	// Check if it has password
	if hash == "" {
		return CredentialDoesNotHavePassword
	}

	// Check on "before" stage
	if err := lookupResult.CheckLogin(types.BeforePasswordCheck); err != nil {
		return err
	}

	// Check the password
	if err := lookupResult.HashingEngine().Validate(password, hash); err != nil {
		return err
	}

	// Check on "after" stage
	if err := lookupResult.CheckLogin(types.AfterPasswordCheck); err != nil {
		return err
	}

	return nil
}


// Sets a new password to the given credential. The given credential must be
// a (*T) value (with T being a struct type).
func SetPassword(credential stub.Credential, password string) error {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(credential) {
		return StructPointerStubExpected
	}

	// Hash the password and store it
	if hash, err := credential.HashingEngine().Hash(password); err != nil {
		return err
	} else {
		credential.SetHashedPassword(hash)
		return nil
	}
}


// Clears the password from the given credential. The given credential must be
// a (*T) value (with T being a struct type).
func ClearPassword(credential stub.Credential) error {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(credential) {
		return StructPointerStubExpected
	}

	// Clear the password
	credential.ClearPassword()
	return nil
}