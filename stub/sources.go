package stub


// There will be one lookup per database engine, most likely.
// The first argument will serve both as a placeholder for the
// result and as specifier of the model and case-sensitivity
// of the search.
type Source interface {
	Lookup(resultHolder Credential, identification string) error
}
