package goutil

import (
	"testing"
)

func TestMap(t *testing.T) {
	m := NormalMap(1000)
	m.Store(1, "a")
	m.Store(2, "b")
	m.Store(3, "c")
	m.Store(4, "d")
	m.Store(5, "e")
	m.Store(6, "f")
	var s = make(map[interface{}]int)
	for i := 10000; i > 0; i-- {
		k, _, _ := m.Random()
		s[k]++
	}
	t.Logf("%#v", s)
}
