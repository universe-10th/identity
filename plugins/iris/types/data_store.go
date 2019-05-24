package types


// A data store interface. Used to get,
// set or delete values.
type DataStore interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Delete(key string) bool
}
