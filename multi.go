package identity

import "github.com/universe-10th/identity/stub"


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