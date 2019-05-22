package identity

import "errors"


var SourceExpected = errors.New("a source is expected")
var StructPointerStubExpected = errors.New("only pointer-kind credential stubs are allowed")
var CredentialDoesNotHavePassword = errors.New("fetched credential does not have password")
var InvalidRealm = errors.New("invalid realm")
var NoMultiMatch = errors.New("no match in any realm")