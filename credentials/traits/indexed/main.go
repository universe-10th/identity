package indexed

// This trait knows its lookup index by
// its own. By contract, must return the
// lookup index a related source would
// use to retrieve this instance.
type Indexed interface {
	Index() interface{}
}
