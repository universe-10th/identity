package types

import (
	"github.com/universe-10th/identity/stub"
	"github.com/kataras/iris/core/errors"
)

type CredentialManager struct {
	clear func() error
	set func(credential stub.Credential) error
	get func() (stub.Credential, error)
}


var NotImplemented error = errors.New("Not implemented")


func NewCredentialManager(
	clear func() error,
	set func(credential stub.Credential) error,
    get func() (stub.Credential, error),
) *CredentialManager {
	return &CredentialManager{clear, set, get}
}


func (credentialManager *CredentialManager) Clear() error {
	if credentialManager.clear == nil {
		return NotImplemented
	} else {
		return credentialManager.clear()
	}
}


func (credentialManager *CredentialManager) Set(credential stub.Credential) error {
	if credentialManager.set == nil {
		return NotImplemented
	} else {
		return credentialManager.set(credential)
	}
}


func (credentialManager *CredentialManager) Get() (stub.Credential, error) {
	if credentialManager.get == nil {
		return nil, NotImplemented
	} else {
		return credentialManager.get()
	}
}