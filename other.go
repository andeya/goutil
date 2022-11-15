package goutil

import (
	"reflect"
	"unsafe"

	"github.com/andeya/ameda"
)

// AddrInt returns a pointer int representing the address of i.
func AddrInt(i int) *int {
	return &i
}

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

// InitPointer initializes null pointer.
func InitPointer(v reflect.Value) bool {
	return ameda.InitPointer(v)
}

// DereferenceType dereference, get the underlying non-pointer type.
func DereferenceType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// DereferenceValue dereference and unpack interface,
// get the underlying non-pointer and non-interface value.
func DereferenceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DereferencePtrValue returns the underlying non-pointer type value.
func DereferencePtrValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// DereferenceIfaceValue returns the value of the underlying type that implements the interface v.
func DereferenceIfaceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DereferenceImplementType returns the underlying type of the value that implements the interface v.
func DereferenceImplementType(v reflect.Value) reflect.Type {
	return DereferenceType(DereferenceIfaceValue(v).Type())
}

// DereferenceSlice convert []*T to []T.
func DereferenceSlice(v reflect.Value) reflect.Value {
	m := v.Len() - 1
	if m < 0 {
		return reflect.New(reflect.SliceOf(DereferenceType(v.Type().Elem()))).Elem()
	}
	s := make([]reflect.Value, m+1)
	for ; m >= 0; m-- {
		s[m] = DereferenceValue(v.Index(m))
	}
	v = reflect.New(reflect.SliceOf(s[0].Type())).Elem()
	v = reflect.Append(v, s...)
	return v
}

// ReferenceSlice convert []T to []*T, the ptrDepth is the count of '*'.
func ReferenceSlice(v reflect.Value, ptrDepth int) reflect.Value {
	if ptrDepth <= 0 {
		return v
	}
	m := v.Len() - 1
	if m < 0 {
		return reflect.New(reflect.SliceOf(ReferenceType(v.Type().Elem(), ptrDepth))).Elem()
	}
	s := make([]reflect.Value, m+1)
	for ; m >= 0; m-- {
		s[m] = ReferenceValue(v.Index(m), ptrDepth)
	}
	v = reflect.New(reflect.SliceOf(s[0].Type())).Elem()
	v = reflect.Append(v, s...)
	return v
}

// ReferenceType convert T to *T, the ptrDepth is the count of '*'.
func ReferenceType(t reflect.Type, ptrDepth int) reflect.Type {
	for ; ptrDepth > 0; ptrDepth-- {
		t = reflect.PtrTo(t)
	}
	return t
}

// ReferenceValue convert T to *T, the ptrDepth is the count of '*'.
func ReferenceValue(v reflect.Value, ptrDepth int) reflect.Value {
	for ; ptrDepth > 0; ptrDepth-- {
		vv := reflect.New(v.Type())
		vv.Elem().Set(v)
		v = vv
	}
	return v
}

// IsLittleEndian determine whether the current system is little endian.
func IsLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return (b == 0x04)
}

const (
	// Is64BitPlatform Whether the current system is a 64-bit platform
	Is64BitPlatform bool = (32 << (^uint(0) >> 63)) == 64
	// Is32BitPlatform Whether the current system is a 32-bit platform
	Is32BitPlatform bool = (32 << (^uint(0) >> 63)) == 0
)
