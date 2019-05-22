package identity

import (
	"github.com/universe-10th/identity/stub"
	"reflect"
)


// A credential manager combines a lookup source and a credential
// prototype to be used inside an internal factory to generate new
// records as needed by marshalling or even logging in.
type CredentialsManager struct {
	source stub.Source
	factoryType reflect.Type
}


// Creates a credentials manager, requiring both the lookup source and
// the prototype.
func NewCredentialsManager(source stub.Source, prototype stub.Credential) (*CredentialsManager, error) {
	// Ensure only a pointer to a struct is used as prototype
	if !prototypeIsAStructPtr(prototype) {
		return nil, StructPointerStubExpected
	}

	// Ensure a source is given
	if source == nil {
		return nil, SourceExpected
	}

	return &CredentialsManager{source, reflect.TypeOf(prototype).Elem()}, nil
}


// Creates an instance of the given type (in heap - returns a pointer) and makes
// an interface out of it.
func (credentialsManager *CredentialsManager) factory() stub.Credential {
	return reflect.New(credentialsManager.factoryType).Interface().(stub.Credential)
}


// Unmarshal a credential: it performs the appropriate lookup given its key.
// Used when de-serializing its ID in a live session.
func (credentialsManager *CredentialsManager) Unmarshal(pk interface{}) (stub.Credential, error) {
	// Find the credential in database, using source and factory.
	holder := credentialsManager.factory()
	if err := credentialsManager.source.ByPrimaryKey(holder, pk); err != nil {
		return nil, err
	} else {
		return holder, nil
	}
}


// Lookup a credential: it performs the appropriate lookup given its identification.
// Used when logging in.
func (credentialsManager *CredentialsManager) Lookup(pk interface{}) (stub.Credential, error) {
	// Find the credential in database, using source and factory.
	holder := credentialsManager.factory()
	if err := credentialsManager.source.ByIdentification(holder, pk); err != nil {
		return nil, err
	} else {
		return holder, nil
	}
}

