package identity

import (
	"github.com/universe-10th/identity/stub"
	"reflect"
	"github.com/universe-10th/identity/support/types"
)


// A (credential) realm combines a lookup source and a credential
// prototype to be used inside an internal factory to generate new
// records as needed by marshalling or even logging in.
type Realm struct {
	source stub.Source
	factoryType reflect.Type
}


// Creates a realm, requiring both the lookup source and the prototype.
func NewRealm(source stub.Source, prototype stub.Credential) (*Realm, error) {
	// Ensure only a pointer to a struct is used as prototype
	if !prototypeIsAStructPtr(prototype) {
		return nil, StructPointerStubExpected
	}

	// Ensure a source is given
	if source == nil {
		return nil, SourceExpected
	}

	return &Realm{source, reflect.TypeOf(prototype).Elem()}, nil
}


// Creates an instance of the given type (in heap - returns a pointer) and makes
// an interface out of it.
func (realm *Realm) factory() stub.Credential {
	return reflect.New(realm.factoryType).Interface().(stub.Credential)
}


// Unmarshal a credential: it performs the appropriate lookup given its key.
// Used when de-serializing its ID in a live session.
func (realm *Realm) Unmarshal(pk interface{}) (stub.Credential, error) {
	// Find the credential in database, using source and factory.
	holder := realm.factory()
	if err := realm.source.ByPrimaryKey(holder, pk); err != nil {
		return nil, err
	} else {
		return holder, nil
	}
}


// Lookup a credential: it performs the appropriate lookup given its identification.
// Used when logging in.
func (realm *Realm) Lookup(pk interface{}) (stub.Credential, error) {
	// Find the credential in database, using source and factory.
	holder := realm.factory()
	if err := realm.source.ByIdentification(holder, pk); err != nil {
		return nil, err
	} else {
		return holder, nil
	}
}


// With realms, now we can implement multi-realms.


// A MultiRealm will work with many credentials managers together.
// You will give a prefix to different backends' marshaled keys to avoid collisions
// you'd have due to lookup eventual priority.
type MultiRealm map[string]*Realm


// Unmarshal a credential: it performs the appropriate lookup given a key and realm key.
// Used when de-serializing its ID in a live session.
func (multiRealm MultiRealm) Unmarshal(realm string, pk interface{}) (stub.Credential, error) {
	if manager, ok := multiRealm[realm]; ok {
		return manager.Unmarshal(pk)
	} else {
		return nil, InvalidRealm
	}
}


// Lookup a credential: it performs the appropriate lookup given an identification and realm.
// Used when logging in (in a specific realm).
func (multiRealm MultiRealm) Lookup(realm string, identification interface{}) (stub.Credential, error) {
	// A single-realm check
	if manager, ok := multiRealm[realm]; ok {
		return manager.Lookup(identification)
	} else {
		return nil, InvalidRealm
	}
}


// Lookup a credential in every realm. The first match will be considered a success, and also
// its realm key will be returned.
func (multiRealm MultiRealm) MultiLookup(identification interface{}) (string, stub.Credential, error) {
	for realm, manager := range multiRealm {
		if credential, err := manager.Lookup(identification); err == nil && credential != nil {
			return realm, credential, err
		}
	}
	return "", nil, NoMultiMatch
}


// Given a realm code, an identification and a password, it tries to perform a login.
// It may fail for several reasons: a database error, an unmatched lookup, a bad password, a
// password-less credential, or another custom login error (after or before the passwords check).
func (multiRealm MultiRealm) Login(realm string, identification interface{}, password string) (string, stub.Credential, error) {
	var credential stub.Credential
	var err error

	// Perform the lookup
	if realm == "" {
		realm, credential, err = multiRealm.MultiLookup(identification)
	} else {
		credential, err = multiRealm.Lookup(realm, identification)
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

	// Succeed: return realm key, credential, and no error
	return realm, credential, nil
}
