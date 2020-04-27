package sources

import "github.com/universe-10th/identity/credentials"

// Retrieves a credential, or an error, by its identifier.
// By contract, and security reasons, if the credential is
// not found, implementors should return login.ErrLoginFailed,
// with a null credential.
type LoginSource interface {
	ByIdentifier(identifier interface{}, template credentials.Credential) (credentials.Credential, error)
}