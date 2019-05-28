package iris

import (
	"github.com/universe-10th/identity"
	"github.com/universe-10th/identity/stub"
	"github.com/universe-10th/identity/support/constants"
)


// A session store interface. Used to get,
// set or delete values. Typically, common
// or JWT sessions will implement this
// interface.
type Session interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Delete(key string) bool
}


// Takes a MultiRealm and interacts with (typically)
// a session. Actually, any session-compatible
// object (satisfying Get, Set, and Delete as common
// sessions do).
//
// This object is actually a helper and you'll rarely
// want to keep it beyond the current request handler.
type WebRealms struct {
	multiRealm identity.MultiRealm
	session    Session
}


// Returns the current credential in session. If any of the
// required data is missing, it performs a logout and returns
// nil. Otherwise, it returns the current credential and realm.
func (webRealms *WebRealms) Current() (string, stub.Credential) {
	realm := webRealms.session.Get(constants.Realm)
	id := webRealms.session.Get(constants.ID)
	if realm == nil || id == nil {
		webRealms.Logout()
		return "", nil
	}

	realmStr := realm.(string)
	credential, err := webRealms.multiRealm.Unmarshal(realmStr, id)
	if credential == nil || err != nil {
		webRealms.Logout()
		return "", nil
	}

	return realmStr, credential
}


// Tries to login a credential given an identification, a (optional,
// by default blank) realm, and a password. If there is a match, that
// data will be inscribed in session: current credential ID and realm.
// Realm and credential will be returned. If an error occurs, then
// the error will be returned instead. Subsequent calls to Current()
// will return this same credential and realm (don't call that in the
// same request: instead preserve these realm/credential values to
// avoid another database call).
func (webRealms *WebRealms) Login(realm string, identification interface{}, password string) (string, stub.Credential, error) {
	realm, credential, err := webRealms.multiRealm.Login(realm, identification, password)
	if err != nil {
		return "", nil, err
	}

	webRealms.session.Set(constants.ID, credential.PrimaryKey())
	webRealms.session.Set(constants.Realm, realm)
	return realm, credential, nil
}


// Performs a logout in the session: it removes the stored ID and
// realm. Subsequent calls to Current() will return nil.
func (webRealms *WebRealms) Logout() {
	webRealms.session.Delete(constants.ID)
	webRealms.session.Delete(constants.Realm)
}


// A factory to create a WebRealm using a given multi-realm reference.
type WebRealmsFactory struct {
	multiRealm identity.MultiRealm
}


// Given a session, creates a *WebRealm instance.
func (factory *WebRealmsFactory) For(session Session) *WebRealms {
	return &WebRealms{factory.multiRealm, session}
}


func Factory(multiRealm identity.MultiRealm) *WebRealmsFactory {
	return &WebRealmsFactory{multiRealm}
}