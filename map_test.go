package goutil

import (
	"sync"
	"testing"
)

func TestRwMap(t *testing.T) {
	m := RwMap(1000)
	m.Store(1, "a")
	m.Store(2, "b")
	m.Store(3, "c")
	m.Store(4, "d")
	m.Store(5, "e")
	m.Store(6, "f")
	t.Logf("Len: %d", m.Len())
	var s = make(map[interface{}]int)
	for i := 10000; i > 0; i-- {
		k, _, _ := m.Random()
		s[k]++
	}
	t.Logf("%#v", s)
}

func TestAtomicMapLen(t *testing.T) {
	m := AtomicMap()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		func(a int) {
			m.Store(a, "a")
			m.LoadOrStore(a, "b")
			m.Len()
			m.Delete(a)
			m.Len()
			m.LoadOrStore(a, "b")
			m.Store(a, "a")
			m.Len()
			wg.Done()
		}(i)
	}

	wg.Wait()

	if a := m.Len(); a != 10 {
		b := m.Len()
		t.Fatalf("Len: expect: 10, but have: %d %d", a, b)
	}
	i := 10 - 1
	for ; i >= 0; i-- {
		m.Delete(i)
	}
	if m.Len() != 0 {
		t.Fatalf("Len: expect: 0, but have: %d", m.Len())
	}
}

func TestAtomicMap(t *testing.T) {
	m := AtomicMap()
	m.Delete(1)
	m.Delete(1)
	m.Store(1, "a")
	m.Store(2, "b")
	m.Store(3, "c")
	m.Store(4, "d")
	m.Store(5, "e")
	m.Store(6, "f")
	m.LoadOrStore(6, "f")
	m.LoadOrStore(6, "f")
	m.Delete(1)
	m.LoadOrStore(1, "a")
	m.Store(1, "a")
	m.Delete(1)
	m.Delete(1)
	if m.Len() != 5 {
		t.Fatalf("Len: expect: 5, but have: %d", m.Len())
	}
	var s = make(map[interface{}]int)
	for i := 10000; i > 0; i-- {
		k, _, _ := m.Random()
		s[k]++
	}
	t.Logf("%#v", s)
}

func TestLoadOrStore(t *testing.T) {
	m := RwMap()
	t.Log(m.LoadOrStore(1, "a"))
	t.Log(m.LoadOrStore(1, "b"))
	t.Log(m.LoadOrStore(1, "c"))
	m = RwMap()
	t.Log(m.LoadOrStore(1, "a"))
	t.Log(m.LoadOrStore(1, "b"))
	t.Log(m.LoadOrStore(1, "c"))
}
