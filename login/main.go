package login

import (
	"errors"
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

// Makes a full login lifecycle function. The returned
// function takes the identification as an arbitrary
// value, the plain-text password as a string, and
// returns either the found and logged credential, or
// an error. To make this function, a login source
// must be used.
func MakeLogin(source sources.LoginSource, template credentials.Credential, steps ...pipeline.PipelineStep) func(interface{}, string) (credentials.Credential, error) {
	return func(identifier interface{}, password string) (credentials.Credential, error) {
		if identifier == nil {
			return nil, ErrNoIdentification
		}

		if credential, err := source.ByIdentifier(identifier); err != nil {
			// These steps are dumb and intended to prevent
			// time correlation attacks to distinguish the
			// case of invalid password and the case of
			// credential not being found.
			for _, step := range steps {
				step.Login(template, password)
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