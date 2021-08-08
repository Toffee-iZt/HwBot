package vktypes

import "reflect"

var reflectTypes = make(map[string]reflect.Type)

// Reg registers new type.
func Reg(s string, i interface{}) {
	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	reflectTypes[s] = typ
}

// Alloc allocates a new zero value of a specific type
// and returns an interface with a pointer to the value.
func Alloc(t string) interface{} {
	typ := reflectTypes[t]
	if typ == nil {
		return nil
	}
	return reflect.New(typ).Interface()
}
