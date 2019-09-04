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
