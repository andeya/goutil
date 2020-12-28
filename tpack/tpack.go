package tpack

import (
	"github.com/henrylee2cn/ameda"
)

// U go underlying type data
type U = ameda.Value

// Unpack unpacks i to go underlying type data.
// Signature:
//  func Unpack(i interface{}) U
var Unpack = ameda.ValueOf

// From gets go underlying type data from reflect.Value.
// Signature:
//  func From(v reflect.Value) U
var From = ameda.ValueFrom

// RuntimeTypeID returns the underlying type ID in current runtime from reflect.Type.
// NOTE:
//  *A and A returns the different runtime type ID;
//  It is 10 times performance of t.String().
// Signature:
//  func RuntimeTypeID(t reflect.Type) uintptr
var RuntimeTypeID = ameda.RuntimeTypeID
