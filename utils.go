package goutil

import (
	"reflect"
)

// InitAndGetString if strPtr is empty string, initialize it with def,
// and return the final value.
func InitAndGetString(strPtr *string, def string) string {
	if strPtr == nil {
		return def
	}
	if *strPtr == "" {
		*strPtr = def
	}
	return *strPtr
}

// DerefType dereference, get the underlying non-pointer type.
func DerefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// DerefValue dereference and unpack interface,
// get the underlying non-pointer and non-interface value.
func DerefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}
