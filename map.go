package goutil

import (
	"sync"
)

// Map is a concurrent map with loads, stores, and deletes.
// It is safe for multiple goroutines to call a Map's methods concurrently.
type Map interface {
	// Load returns the value stored in the map for a key, or nil if no
	// value is present.
	// The ok result indicates whether value was found in the map.
	Load(key interface{}) (value interface{}, ok bool)
	// Store sets the value for a key.
	Store(key, value interface{})
	// LoadOrStore returns the existing value for the key if present.
	// Otherwise, it stores and returns the given value.
	// The loaded result is true if the value was loaded, false if stored.
	LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
	// Delete deletes the value for a key.
	Delete(key interface{})
	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, range stops the iteration.
	Range(f func(key, value interface{}) bool)
	// InexactLen returns the length of the map.
	// Note:
	//  the count implemented using sync.Map may be inaccurate;
	//  the count implemented using NormalMap is accurate.
	InexactLen() int
}

// NormalMap make a new concurrent safe map with sync.RWRWMutex.
// The normal Map is high-performance mapping under low concurrency conditions.
func NormalMap(capacity ...int) Map {
	var cap int
	if len(capacity) > 0 {
		cap = capacity[0]
	}
	return &normalMap{
		data: make(map[interface{}]interface{}, cap),
	}
}

// normalMap concurrent secure data storage,
// which is high-performance mapping under low concurrency conditions.
type normalMap struct {
	data map[interface{}]interface{}
	rwmu sync.RWMutex
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *normalMap) Load(key interface{}) (value interface{}, ok bool) {
	m.rwmu.RLock()
	value, ok = m.data[key]
	m.rwmu.RUnlock()
	return value, ok
}

// Store sets the value for a key.
func (m *normalMap) Store(key, value interface{}) {
	m.rwmu.Lock()
	m.data[key] = value
	m.rwmu.Unlock()
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *normalMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	m.rwmu.Lock()
	actual, loaded = m.data[key]
	m.data[key] = value
	if !loaded {
		actual = value
	}
	m.rwmu.Unlock()
	return actual, loaded
}

// Delete deletes the value for a key.
func (m *normalMap) Delete(key interface{}) {
	m.rwmu.Lock()
	delete(m.data, key)
	m.rwmu.Unlock()
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (m *normalMap) Range(f func(key, value interface{}) bool) {
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// InexactLen returns the length of the map.
// Note: the count is accurate.
func (m *normalMap) InexactLen() int {
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	return len(m.data)
}
