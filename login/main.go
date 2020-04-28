package login

import (
	"errors"
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/login/pipeline"
)

// Generic error to return in most of the pipeline
// steps and when a source cannot find a credential
// by its identification.
var ErrLoginFailed = errors.New("login failed")

// Error to return whether a nil identification was
// provided.
var ErrNoIdentification = errors.New("no identification provided")

// A login realm is a class combining a full pipeline
// and a source. It only provides one method: Login,
// which takes the identifier and password to attempt
// a user lookup and then the actual login process by
// running all the elements in the pipe.
type LoginRealm struct {
	source credentials.Source
	steps  []pipeline.PipelineStep
}

// Makes a full login lifecycle function. The returned
// function takes the identification as an arbitrary
// value, the plain-text password as a string, and
// returns either the found and logged credential, or
// an error. To make this function, a login source
// must be used. A template credential is used to both
// serve as factory and dummy.
func (realm *LoginRealm) Login(identifier interface{}, password string) (credentials.Credential, error) {
	if identifier == nil {
		return nil, ErrNoIdentification
	}

	if credential, err := realm.source.ByIdentifier(identifier); err != nil {
		// These steps are dumb and intended to prevent
		// time correlation attacks to distinguish the
		// case of invalid password and the case of
		// credential not being found.
		dummy := realm.source.Dummy()
		for _, step := range realm.steps {
			_ = step.Login(dummy, password)
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
