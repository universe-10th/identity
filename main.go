package identity

import "reflect"


func prototypeIsAStructPtr(prototype interface{}) bool {
	// nil interface is not allowed
	if prototype == nil {
		return false
	}
	// Only prototypes that are a pointer types are allowed.
	prType := reflect.TypeOf(prototype)
	if prType.Kind() != reflect.Ptr {
		return false
	}
	// Also the indirect type must be a struct.
	return prType.Elem().Kind() != reflect.Struct
}