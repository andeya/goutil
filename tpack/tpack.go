package tpack

import (
	"reflect"
	"unsafe"
)

// T go underlying type data
type T struct {
	typeID int32
	i      interface{}
}

// Unpack unpack i to go underlying type data.
func Unpack(i interface{}) T {
	p := *(*uintptr)(unsafe.Pointer(&i))
	v := T{
		typeID: *(*int32)((unsafe.Pointer(p + rtypeStrOffset))),
		i:      i,
	}
	return v
}

// TypeID returns the underlying type ID.
// It is 60 times performance of reflect.TypeOf(i).String()
func (t T) TypeID() int32 {
	return t.typeID
}

// TypeOf is equivalent to reflect.TypeOf.
func (t T) TypeOf() reflect.Type {
	return reflect.TypeOf(t.i)
}

// ValueOf is equivalent to reflect.ValueOf.
func (t T) ValueOf() reflect.Value {
	return reflect.ValueOf(t.i)
}

// TypeID get underlying type ID from reflect.Type.
// It is 60 times performance of t.String()
func TypeID(t reflect.Type) int32 {
	ptr := elemUintptr(uintptr(unsafe.Pointer(&t)) + ptrOffset)
	return *(*int32)((unsafe.Pointer(ptr + rtypeStrOffset)))
}

func elemUintptr(ptr uintptr) uintptr {
	return *(*uintptr)(unsafe.Pointer(ptr))
}

// func getTypeIDReference(i interface{}) int32 {
// 	return int32((*(*reflectValue)(unsafe.Pointer(&i))).typ.str)
// }

var (
	rtypeStrOffset = func() uintptr {
		return unsafe.Offsetof(e.typ.str)
	}()
	ptrOffset = func() uintptr {
		return unsafe.Offsetof(e.ptr)
	}()
	e = reflectValue{typ: new(rtype)}
)

type tflag uint8

type (
	reflectValue struct {
		typ *rtype
		ptr unsafe.Pointer
		flag
	}
	rtype struct {
		size       uintptr
		ptrdata    uintptr  // number of bytes in the type that can contain pointers
		hash       uint32   // hash of type; avoids computation in hash tables
		tflag      tflag    // extra type information flags
		align      uint8    // alignment of variable with this type
		fieldAlign uint8    // alignment of struct field with this type
		kind       uint8    // enumeration for C
		alg        *typeAlg // algorithm table
		gcdata     *byte    // garbage collection data
		str        nameOff  // string form
		ptrToThis  typeOff  // type for pointer to this type, may be zero
	}
	typeAlg struct {
		hash  func(unsafe.Pointer, uintptr) uintptr
		equal func(unsafe.Pointer, unsafe.Pointer) bool
	}
	nameOff int32 // offset to a name
	typeOff int32 // offset to an *rtype
	flag    uintptr
)
