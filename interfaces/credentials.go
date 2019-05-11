package interfaces

import "github.com/luismasuelli/go-identity/types"


/**
 * Credential interface will be implemented on particular
 *   models we're interested about to be credentials (one
 *   example: users).
 */
type Credential interface {
	// Tells which field has the identification.
	// This, instead of retrieving the identification:
	//   It will allow us to use that field to perform
	//   a database search.
	IdentificationField() string
	// Tells whether the identification lookup is case
	//   sensitive or not.
	IdentificationIsCaseSensitive() bool
	// A mean to set the password. It will most likely
	//   store the (hashed) password inside a particular
	//   field in the credential.
	SetHashedPassword(string)
	// A mean to get the hashed password.
	HashedPassword() string
	// A mean to clear the password. Credentials with no
	//   password will not authenticate at all.
	ClearPassword()
	// A reference to the hashing engine to use.
	HashingEngine() PasswordHashingEngine
	// Tells whether this credential can login or not.
	// This method should not consider whether the
	//   credential has password or the passwords match:
	//   such check will also run in a different moment.
	CheckLogin(stage types.LoginStage) error
}
