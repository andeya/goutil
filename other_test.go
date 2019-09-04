package goutil

import (
	"reflect"
	"testing"
)

func TestReferenceSlice(t *testing.T) {
	v := reflect.ValueOf([]int{1, 2})
	v = ReferenceSlice(v, 1)
	ret := v.Interface().([]*int)
	t.Log(ret)
}

func TestDereferenceSlice(t *testing.T) {
	one := 1
	two := 2
	v := reflect.ValueOf([]*int{&one, &two})
	v = DereferenceSlice(v)
	ret := v.Interface().([]int)
	t.Log(ret)
}
