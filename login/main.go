package login

import (
	"errors"
	"reflect"
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/login/sources"
	"github.com/universe-10th/identity/login/pipeline"
)

// Generic error to return in most of the pipeline
// steps and when a source cannot find a credential
// by its identification.
var ErrLoginFailed = errors.New("login failed")

// Error to return whether a nil identification was
// provided.
var ErrNoIdentification = errors.New("no identification provided")

// A login realm is a function taking identifier and
// password, and returning a credential or an error
// after attempting a login.
type LoginRealm func(interface{}, string) (credentials.Credential, error)

// Makes a full login lifecycle function. The returned
// function takes the identification as an arbitrary
// value, the plain-text password as a string, and
// returns either the found and logged credential, or
// an error. To make this function, a login source
// must be used. A template credential is used to both
// serve as factory and dummy.
func MakeLoginRealm(source sources.LoginSource, template credentials.Credential, steps ...pipeline.PipelineStep) LoginRealm {
	credType := reflect.TypeOf(template)
	var factory func() credentials.Credential
	if credType.Kind() == reflect.Ptr {
		credElemType := credType.Elem()
		factory = func() credentials.Credential {
			return reflect.New(credElemType).Interface().(credentials.Credential)
		}
	} else {
		factory = func() credentials.Credential {
			return reflect.New(credType).Elem().Interface().(credentials.Credential)
		}
	}

	return func(identifier interface{}, password string) (credentials.Credential, error) {
		if identifier == nil {
			return nil, ErrNoIdentification
		}

		if credential, err := source.ByIdentifier(identifier, template); err != nil {
			// These steps are dumb and intended to prevent
			// time correlation attacks to distinguish the
			// case of invalid password and the case of
			// credential not being found.
			dummy := factory()
			for _, step := range steps {
				step.Login(dummy, password)
			}
			return nil, err
		} else {
			for _, step := range steps {
				if stepErr := step.Login(credential, password); stepErr != nil {
					return nil, stepErr
				}
			}
			return credential, nil
		}
	}
}