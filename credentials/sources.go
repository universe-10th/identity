package credentials

import (
	"errors"
	"reflect"
)

// Brokers perform load and save operations over an existing
// credential (they are not intended to create a new credential).
// The login operation will involve calling the ByIdentifier
// method, while the other two are intended for the different
// password change pipelines. Sources will also be able to tell
// whether they allow, or not, a given implementation of the
// Credential interface. This usually involves other traits
// being implemented in the Credential, and/or things like the
// credential being a struct-type value. Notes: ByIdentifier
// must return (nil, nil) if no credential can be found by its
// identification.
type Broker interface {
	Allows(template Credential) bool
	ByIdentifier(identifier interface{}, template Credential) (Credential, error)
	ByIndex(index interface{}, template Credential) (Credential, error)
	Save(credential Credential) error
}

// Panicked error when attempting to create a source with a
// nil broker instead of an instance.
var ErrNilBroker = errors.New("the given broker is nil")

// Panicked error when attempting to create a source with a
// nil template credential, or a template not allowed by the
// chosen broker.
var ErrBadTemplate = errors.New("the given template is nil or not allowed by the broker")

var ErrNilValueOnSave = errors.New("the credential being saved is nil")

// Returned error when the credential being saved is of a
// different type than the one of the source.
var ErrBadTypeOnSave = errors.New("the credential being saved is of a different type")

// Sources are a combination of an existing broker instance and
// a non-nil Credential instance that will serve as template.
// Sources will proxy the calls to a broker, and also will be
// able to instantiate dummy objects of the same type of the
// given template.
type Source struct {
	broker   Broker
	template Credential
	tmplType reflect.Type
	factory  func() Credential
}

// Creates a new source for a given broker and template, if they
// match together. Panics if either is nil, or the broker does
// not allow it.
func NewSource(broker Broker, template Credential) *Source {
	if broker == nil {
		panic(ErrNilBroker)
	} else if template == nil || !broker.Allows(template) {
		panic(ErrBadTemplate)
	}

	credType := reflect.TypeOf(template)
	var factory func() Credential
	if credType.Kind() == reflect.Ptr {
		credElemType := credType.Elem()
		factory = func() Credential {
			return reflect.New(credElemType).Interface().(Credential)
		}
	} else {
		factory = func() Credential {
			return reflect.New(credType).Elem().Interface().(Credential)
		}
	}
	return &Source{broker, template, credType, factory}
}

// Bypasses its implementation to the broker but using the chosen
// template instance.
func (source *Source) ByIdentifier(identifier interface{}) (Credential, error) {
	return source.broker.ByIdentifier(identifier, source.template)
}

// Bypasses its implementation to the broker but using the chosen
// template instance.
func (source *Source) ByIndex(index interface{}) (Credential, error) {
	return source.broker.ByIndex(index, source.template)
}

// Bypasses its implementation to the broker but using the chosen
// template instance, adding a check type.
func (source *Source) Save(credential Credential) error {
	if credential == nil {
		return ErrNilValueOnSave
	} else if reflect.TypeOf(credential) != source.tmplType {
		return ErrBadTypeOnSave
	} else {
		return source.broker.Save(credential)
	}
}

// Creates a dummy credential object, used for security
// purposes following a fake login cycle.
func (source *Source) Dummy() Credential {
	return source.factory()
}
