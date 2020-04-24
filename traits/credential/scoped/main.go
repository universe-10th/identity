package scoped

// A scope is a particular keyed token that
// is meant to be used as a mean to satisfy
// a requirement (and also to specify an
// atomic requirement).
type Scope interface {
	// The fully-qualified scope key.
	Key() string
	// The scope friendly name.
	Name() string
	// The scope description (optional).
	Description() string
}

// This trait makes the credential hold a
// bunch of scopes that will be used to test
// against scope requirements, which may be
// simple or composed.
type Scoped interface {
	Scopes() map[Scope]bool
}
