package identity

import (
	"github.com/universe-10th/identity/stub"
	"github.com/universe-10th/identity/support/types"
)

// Given a realm, a credentials multi-manager, an identification and a password, it tries to
// perform a login. It may fail for several reasons: a database error, an unmatched lookup, a
// bad password, a password-less credential, or another custom login error (after or before the
// passwords check).
func Login(realm string, managers CredentialsMultiManager, identification interface{}, password string) (string, stub.Credential, error) {
	var credential stub.Credential
	var err error

	// Perform the lookup
	if realm == "" {
		realm, credential, err = managers.MultiLookup(identification)
	} else {
		credential, err = managers.Lookup(realm, identification)
	}
	if err != nil {
		return "", nil, err
	}

	// Check if it has password
	hash := credential.HashedPassword()
	if hash == "" {
		return "", nil, CredentialDoesNotHavePassword
	}

	// Check on "before" stage
	if err := credential.CheckLogin(types.BeforePasswordCheck); err != nil {
		return "", nil, err
	}

	// Check the password
	if err := credential.HashingEngine().Validate(password, hash); err != nil {
		return "", nil, err
	}

	// Check on "after" stage
	if err := credential.CheckLogin(types.AfterPasswordCheck); err != nil {
		return "", nil, err
	}

	// Succeed: return realm, credential, and no error
	return realm, credential, nil
}