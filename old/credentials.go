package old

import "github.com/universe-10th/identity/support/types"

// Credential interface will be implemented on particular
// models we're interested about to be credentials (one
// example: users).
type Credential interface {
	// Retrieves the credential's primary key, used
	// to serialize / store in sessions or JWT.
	PrimaryKey() interface{}
	// Tells which field has the primary key.
	// It will allow us to use that field to perform
	// a database search. It must be consistent with
	// the implementation of PrimaryKey() method.
	PrimaryKeyField() string
	// Tells which field has the identification.
	// It will allow us to use that field to perform
	// a database search.
	IdentificationField() string
	// Tells whether the identification lookup is case
	// sensitive or not (only useful to string fields).
	IdentificationIsCaseSensitive() bool
	// Gets the identification (it must make sense with
	// the value in IdentificationField).
	Identification() interface{}
	// Sets the identification (it must make sense with
	// the value in IdentificationField).
	SetIdentification(identification interface{})
	// A mean to set the password. It will most likely
	// store the (hashed) password inside a particular
	// field in the credential.
	SetHashedPassword(string)
	// A mean to get the hashed password.
	HashedPassword() string
	// A mean to clear the password. Credentials with no
	// password will not authenticate at all.
	ClearPassword()
	// A reference to the hashing engine to use.
	HashingEngine() PasswordHashingEngine
	// Tells whether this credential can login or not.
	// This method should not consider whether the
	// credential has password or the passwords match:
	// such check will also run in a different moment.
	CheckLogin(stage types.LoginStage) error
}

// Super user interface provides a way to tell whether the
// underlying object should be understood as a superuser.
type WithSuperUserFlag interface {
	IsSuperUser() bool
}