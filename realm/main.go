package realm

import (
	"errors"
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/credentials/traits/recoverable"
	"github.com/universe-10th/identity/realm/login"
	"time"
)

// Generic error to return in most of the pipeline
// steps and when a source cannot find a credential
// by its identification.
var ErrLoginFailed = errors.New("login failed")

// Error to return when current password does not match.
var ErrBadCurrentPassword = errors.New("invalid current password")

// Error to return when attempting to use a credential
// that is not recoverable on password reset features.
var ErrNotRecoverable = errors.New("the credential is not a recoverable type")

// Error to return when an invalid token was given to a
// password reset attempt.
var ErrBadToken = errors.New("invalid token on password reset confirm, or password reset was not issued")

// Panicked when a nil source is given to a realm.
var ErrNilSource = errors.New("source is nil")

// Panicked when a nil pipeline step is given to a realm.
var ErrNilPipelineStep = errors.New("pipeline step is nil")

// A login realm is a class combining a full pipeline
// and a source. It only provides one method: Login,
// which takes the identifier and password to attempt
// a user lookup and then the actual login process by
// running all the elements in the pipe.
type Realm struct {
	source *credentials.Source
	steps  []login.PipelineStep
}

// Retrieves a credential by its identifier. This call is directly bypassed to the source.
func (realm *Realm) ByIdentifier(identifier interface{}) (credentials.Credential, error) {
	return realm.source.ByIdentifier(identifier)
}

// Retrieves a credential by its index / key. This call is directly bypassed to the source.
func (realm *Realm) ByIndex(index interface{}) (credentials.Credential, error) {
	return realm.source.ByIndex(index)
}

// Makes a full login lifecycle function. The returned
// function takes the identification as an arbitrary
// value, the plain-text password as a string, and
// returns either the found and logged credential, or
// an error. To make this function, a login source
// must be used. A template credential is used to both
// serve as factory and dummy.
func (realm *Realm) Login(identifier interface{}, password string) (credentials.Credential, error) {
	if credential, err := realm.ByIdentifier(identifier); credential == nil {
		// These steps are dumb and intended to prevent
		// time correlation attacks to distinguish the
		// case of invalid password and the case of
		// credential not being found.
		dummy := realm.source.Dummy()
		for _, step := range realm.steps {
			_ = step.Login(dummy, password)
		}
		// When both credential and error are nil, the
		// ErrLoginFailed will be used instead.
		if err == nil {
			err = ErrLoginFailed
		}
		return nil, err
	} else {
		for _, step := range realm.steps {
			if stepErr := step.Login(credential, password); stepErr != nil {
				return nil, stepErr
			}
		}
		return credential, nil
	}
}

// Attempts a password change, which involves invoking the appropriate hashing.
// The credential will be saved after that.
func (realm *Realm) SetPassword(credential credentials.Credential, password string) error {
	if hashedPassword, err := credential.Hasher().Hash(password); err != nil {
		return err
	} else {
		credential.SetHashedPassword(hashedPassword)
		return realm.source.Save(credential)
	}
}

// Attempts a password unset, which involves deleting the hashed password.
// The credential will be saved after that.
func (realm *Realm) UnsetPassword(credential credentials.Credential) error {
	credential.SetHashedPassword("")
	return realm.source.Save(credential)
}

// Attempts a by-user password change, which involves invoking the appropriate
// hashing and also validating the current password. The credential will be
// saved after that.
func (realm *Realm) ChangePassword(credential credentials.Credential, currentPassword, newPassword string) error {
	if err := credential.Hasher().Validate(currentPassword, credential.HashedPassword()); err != nil {
		return ErrBadCurrentPassword
	} else {
		return realm.SetPassword(credential, newPassword)
	}
}

// Attempts an external, non-logged and to-be-confirmed attempt to reset a password.
// It will set the recovery token and save the credential. This call is only allowed
// if the credential is of a recoverable type.
func (realm *Realm) PreparePasswordReset(credential credentials.Credential, token string, duration time.Duration) error {
	if recoverableCred, ok := credential.(recoverable.Recoverable); !ok {
		return ErrNotRecoverable
	} else {
		recoverableCred.SetRecoveryToken(token, duration)
		return realm.source.Save(credential)
	}
}

// Clears an external, non-logged and to-be-confirmed attempt to reset a password.
// This call is only allowed if the credential is of a recoverable type.
func (realm *Realm) CancelPasswordReset(credential credentials.Credential) error {
	return realm.PreparePasswordReset(credential, "", time.Duration(0))
}

// Confirms an external, non-logged and to-be-confirmed attempt to reset a password.
// This call is only allowed if the credential is of a recoverable type.
func (realm *Realm) ConfirmPasswordReset(credential credentials.Credential, token, password string) error {
	if recoverableCred, ok := credential.(recoverable.Recoverable); !ok {
		return ErrNotRecoverable
	} else if token != recoverableCred.RecoveryToken() || token == "" {
		return ErrBadToken
	} else if hashed, err := credential.Hasher().Hash(password); err != nil {
		return err
	} else {
		credential.SetHashedPassword(hashed)
		recoverableCred.SetRecoveryToken("", time.Duration(0))
		return realm.source.Save(credential)
	}
}

// Creates a new realm.
func NewRealm(source *credentials.Source, steps ...login.PipelineStep) *Realm {
	if source == nil {
		panic(ErrNilSource)
	}

	for _, step := range steps {
		if step == nil {
			panic(ErrNilPipelineStep)
		}
	}

	return &Realm{source: source, steps: steps}
}
