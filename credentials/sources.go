package credentials

// Sources perform load and save operations over an existing
// credential (they are not intended to create a new credential).
// The login operation will involve calling the ByIdentifier
// method, while the other two are intended for the different
// password change pipelines.
type Source interface {
	ByIdentifier(identifier interface{}, template Credential) (Credential, error)
	ByIndex(index interface{}, template Credential) (Credential, error)
	Save(credential Credential)
}
