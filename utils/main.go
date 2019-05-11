package utils

import (
	"fmt"
	"errors"
	"reflect"
	"github.com/jinzhu/gorm"
	"github.com/luismasuelli/gormid/interfaces"
	"github.com/luismasuelli/gormid/types"
)

var StructPointerPrototypeExpected = errors.New("only pointer-kind prototypes are allowed")
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


/**
 * Given a database query, a model prototype, an identification and a password, it tries to
 *   perform a login. It may fail for several reasons: the prototype is not a (*T) type (with
 *   T being a struct type), a database error, a bad password, a password-less credential, or
 *   another custom login error (after or before the passwords check).
 */
func Login(db *gorm.DB, prototype interfaces.Credential, identification string, password string) (interfaces.Credential, error) {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(prototype) {
		return nil, StructPointerPrototypeExpected
	}

	// Find the credential in database
	caseSensitive := prototype.IdentificationIsCaseSensitive()
	lookupResult := reflect.New(reflect.ValueOf(prototype).Elem().Type()).Interface().(interfaces.Credential)
	query := ""
	if caseSensitive {
		query = fmt.Sprintf("%s = ?", prototype.IdentificationField())
	} else {
		query = fmt.Sprintf("UPPER(%s) = UPPER(?)", prototype.IdentificationField())
	}
	if err := db.Where(query, identification).First(lookupResult).Error; err != nil {
		return nil, err
	}

	hash := lookupResult.HashedPassword()

	// Check if it has password
	if hash == "" {
		return nil, CredentialDoesNotHavePassword
	}

	// Check on "before" stage
	if err := lookupResult.CheckLogin(types.BeforePasswordCheck); err != nil {
		return nil, err
	}

	// Check the password
	if err := lookupResult.HashingEngine().Validate(password, hash); err != nil {
		return nil, err
	}

	// Check on "after" stage
	if err := lookupResult.CheckLogin(types.AfterPasswordCheck); err != nil {
		return nil, err
	}

	return lookupResult, nil
}


/**
 * Sets a new password to the given credential. The given credential must be
 *   a (*T) value (with T being a struct type).
 */
func SetPassword(credential interfaces.Credential, password string) error {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(credential) {
		return StructPointerPrototypeExpected
	}

	// Hash the password and store it
	if hash, err := credential.HashingEngine().Hash(password); err != nil {
		return err
	} else {
		credential.SetHashedPassword(hash)
		return nil
	}
}


/**
 * Clears the password from the given credential. The given credential must be
 *   a (*T) value (with T being a struct type).
 */
func ClearPassword(credential interfaces.Credential) error {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(credential) {
		return StructPointerPrototypeExpected
	}

	// Clear the password
	credential.ClearPassword()
	return nil
}