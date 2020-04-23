package identity

// Sets a new password to the given credential. The given credential must be
// a (*T) value (with T being a struct type).
func SetPassword(credential Credential, password string) error {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(credential) {
		return StructPointerStubExpected
	}

	// Hash the password and store it
	if hash, err := credential.HashingEngine().Hash(password); err != nil {
		return err
	} else {
		credential.SetHashedPassword(hash)
		return nil
	}
}

// Clears the password from the given credential. The given credential must be
// a (*T) value (with T being a struct type).
func ClearPassword(credential Credential) error {
	// Ensure only a pointer to a struct enters here
	if !prototypeIsAStructPtr(credential) {
		return StructPointerStubExpected
	}

	// Clear the password
	credential.ClearPassword()
	return nil
}
