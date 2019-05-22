package utils

import (
	"github.com/kataras/iris"
	"github.com/universe-10th/identity/support/constants"
	"github.com/universe-10th/identity/stub"
)


const credentialSetKey = constants.Namespace + ":credential:set"
const credentialGetKey = constants.Namespace + ":credential:get"


// Users that make use of this IRIS plugin will make use of this
// function to get the credential in the current context.
func SetCredential(context iris.Context, credential stub.Credential) {
	context.Values().Set(credentialSetKey, credential)
}


// Users that make use of this IRIS plugin will make use of this
// function to get the credential in the current context.
func GetCredential(context iris.Context) stub.Credential {
	value, _ := context.Values().Get(credentialGetKey).(stub.Credential)
	return value
}


// Intended to be used by IRIS-compatible middleware functions.
// This one will be used (if a user is available according to the
// middleware function) to give a current credential to the context.
func ProvideCredential(context iris.Context, credential stub.Credential) {
	context.Values().Set(credentialGetKey, credential)
}


// Intended to be used by IRIS-compatible middleware functions.
// This one will be used to retrieve the credential which was set
// by the user in any handler.
func ReceiveCredential(context iris.Context) stub.Credential {
	value, _ := context.Values().Get(credentialSetKey).(stub.Credential)
	return value
}