package traits


// This trait knows its identification
// by its own. By contract, must return
// the identification a related source
// would use to retrieve this instance.
type SelfIdentified interface {
	Identification() interface{}
}
