package login

import (
	"errors"
	"github.com/universe-10th/identity/credentials"
)

var ErrLoginFailed = errors.New("login failed")

// A login pipeline step performs a check on a given
// credential (this implies the credential exists and
// was successfully retrieved) to tell whether it may
// login or not. The attempted password comes as the
// second argument. The result value must be nil if
// the step approves the login attempt, and an error
// instance if it rejects the login attempt.
type PipelineStep interface {
	Login(credential credentials.Credential, password string) error
}


// TODO create the actual pipeline (it will need the realm instance).