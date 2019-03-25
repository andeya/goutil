package tpack

import (
	"reflect"
	"unsafe"
)

// T go underlying type data
type T struct {
	typPtr uintptr
	ptr    unsafe.Pointer
	_iPtr  unsafe.Pointer // avoid being GC
}

// Unpack unpacks i to go underlying type data.
func Unpack(i interface{}) T {
	return newT(unsafe.Pointer(&i))
}

// From gets go underlying type data from reflect.Value.
func From(v reflect.Value) T {
	return newT(unsafe.Pointer(&v))
}

func newT(iPtr unsafe.Pointer) T {
	return T{
		typPtr: *(*uintptr)(iPtr),
		ptr:    pointerElem(unsafe.Pointer(uintptr(iPtr) + ptrOffset)),
		_iPtr:  iPtr,
	}
}

// RuntimeTypeID returns the underlying type ID in current runtime from reflect.Type.
// NOTE:
//  *A and A returns the same runtime type ID;
//  It is 10 times performance of t.String().
func RuntimeTypeID(t reflect.Type) int32 {
	typPtr := uintptrElem(uintptr(unsafe.Pointer(&t)) + ptrOffset)
	return *(*int32)(unsafe.Pointer(typPtr + rtypeStrOffset))
}

// RuntimeTypeID gets the underlying type ID in current runtime.
// NOTE:
//  *A and A gets the same runtime type ID;
//  It is 10 times performance of reflect.TypeOf(i).String().
func (t T) RuntimeTypeID() int32 {
	return *(*int32)(unsafe.Pointer(t.typPtr + rtypeStrOffset))
}

// Kind gets the reflect.Kind fastly.
func (t T) Kind() reflect.Kind {
	return kind(t.typPtr)
}

// UnderlyingKind gets the underlying reflect.Kind fastly.
func (t T) UnderlyingKind() reflect.Kind {
	k := t.Kind()
	typPtr := t.typPtr
	for k == reflect.Ptr || k == reflect.Interface {
		k, typPtr = underlying(k, typPtr)
	}
	return k
}

// Elem returns the value T that the interface i contains
// or that the pointer i points to.
func (t T) Elem() T {
	k := t.Kind()
	if k != reflect.Ptr && k != reflect.Interface {
		return t
	}
	k, t.typPtr = underlying(k, t.typPtr)
	if k == reflect.Ptr || k == reflect.Interface {
		t.ptr = pointerElem(t.ptr)
	}
	return t
}

// UnderlyingElem returns the underlying value T that the interface i contains
// or that the pointer i points to.
func (t T) UnderlyingElem() T {
	for r := t.Elem(); r != t; r = t.Elem() {
		t = r
	}
	return t
}

// Pointer gets the pointer of i.
// NOTE:
//  *A and A, gets diffrent pointer
func (t T) Pointer() uintptr {
	k := t.Kind()
	switch k {
	case reflect.Invalid:
		return 0
	case reflect.Func:
		return uintptr(*(*unsafe.Pointer)(t.ptr))
	case reflect.Slice:
		return uintptrElem(uintptr(t.ptr)) + sliceDataOffset
	default:
		return uintptr(t.ptr)
	}
}

func underlying(k reflect.Kind, typPtr uintptr) (reflect.Kind, uintptr) {
	typPtr = uintptrElem(typPtr + elemOffset)
	return kind(typPtr), typPtr
}

func kind(typPtr uintptr) reflect.Kind {
	k := *(*uint8)(unsafe.Pointer(typPtr + kindOffset))
	return reflect.Kind(k & kindMask)
}

func uintptrElem(ptr uintptr) uintptr {
	return *(*uintptr)(unsafe.Pointer(ptr))
}

func pointerElem(p unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(p)
}

var (
	ptrOffset = func() uintptr {
		return unsafe.Offsetof(e.word)
	}()
	rtypeStrOffset = func() uintptr {
		return unsafe.Offsetof(e.typ.str)
	}()
	kindOffset = func() uintptr {
		return unsafe.Offsetof(e.typ.kind)
	}()
	elemOffset = func() uintptr {
		return unsafe.Offsetof(new(ptrType).elem)
	}()
	sliceLenOffset = func() uintptr {
		return unsafe.Offsetof(new(reflect.SliceHeader).Len)
	}()
	sliceDataOffset = func() uintptr {
		return unsafe.Offsetof(new(reflect.SliceHeader).Data)
	}()
	e = emptyInterface{typ: new(rtype)}
)

const (
	kindMask = (1 << 5) - 1
)

type (
	// reflectValue struct {
	// 	typ *rtype
	// 	ptr unsafe.Pointer
	// 	flag
	// }
	emptyInterface struct {
		typ  *rtype
		word unsafe.Pointer
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
	ptrType struct {
		rtype
		elem *rtype // pointer element (pointed at) type
	}
	typeAlg struct {
		hash  func(unsafe.Pointer, uintptr) uintptr
		equal func(unsafe.Pointer, unsafe.Pointer) bool
	}
	nameOff int32 // offset to a name
	typeOff int32 // offset to an *rtype
	flag    uintptr
	tflag   uint8
)
