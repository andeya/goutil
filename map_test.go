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
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("after clear len: %d", m.Len())
	}
}

func TestLoadOrStore(t *testing.T) {
	m := RwMap()
	if v, loaded := m.LoadOrStore(1, "rw-a"); v != "rw-a" || loaded {
		t.Fatalf("v: %v, loaded: %v", v, loaded)
	}
	if v, loaded := m.LoadOrStore(1, "rw-b"); v != "rw-a" || !loaded {
		t.Fatalf("v: %v, loaded: %v", v, loaded)
	}
	if v, loaded := m.LoadOrStore(1, "rw-c"); v != "rw-a" || !loaded {
		t.Fatalf("v: %v, loaded: %v", v, loaded)
	}
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("after clear len: %d", m.Len())
	}
	m = AtomicMap()
	if v, loaded := m.LoadOrStore(1, "atomic-a"); v != "atomic-a" || loaded {
		t.Fatalf("v: %v, loaded: %v", v, loaded)
	}
	if v, loaded := m.LoadOrStore(1, "atomic-b"); v != "atomic-a" || !loaded {
		t.Fatalf("v: %v, loaded: %v", v, loaded)
	}
	if v, loaded := m.LoadOrStore(1, "atomic-c"); v != "atomic-a" || !loaded {
		t.Fatalf("v: %v, loaded: %v", v, loaded)
	}
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("after clear len: %d", m.Len())
	}
}

func TestAtomicMap(t *testing.T) {
	m := AtomicMap()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(a int) {
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

	var s = make(map[interface{}]int)
	for i := 100000; i > 0; i-- {
		k, _, _ := m.Random()
		s[k]++
	}
	t.Logf("%#v", s)

	i := 10 - 1
	for ; i >= 0; i-- {
		m.Delete(i)
	}
	if m.Len() != 0 {
		t.Fatalf("Len: expect: 0, but have: %d", m.Len())
	}
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("after clear len: %d", m.Len())
	}
}

func TestLen(t *testing.T) {
	m := AtomicMap()
	var wg sync.WaitGroup
	for i := 1; i <= 100000; i++ {
		wg.Add(1)
		go func(a int) {
			m.Store(a, "a")
			wg.Done()
		}(i)
	}
	m.Range(func(k, v interface{}) bool {
		m.Delete(k)
		return true
	})
	wg.Wait()
	m.Range(func(k, v interface{}) bool {
		m.Delete(k)
		return true
	})
	if a := m.Len(); a != 0 {
		t.Fatalf("Len: expect: 0, but have: %d", a)
	}
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("after clear len: %d", m.Len())
	}
}

func TestOther(t *testing.T) {
	m := RwMap()
	t.Log(m.Load("key"))
	m.Range(func(_, _ interface{}) bool {
		return true
	})
	m.Len()
	m.Random()

	m = AtomicMap()
	t.Log(m.Load("key"))
	m.Range(func(_, _ interface{}) bool {
		return true
	})
	m.Len()
	m.Random()
}
