package connpool

import (
	"sync"
	"sync/atomic"
)

// ConnPools stores ConnPool
type ConnPools struct {
	// stores 'map[string]ConnPool',
	// one server node has one connection pool.
	pools atomic.Value
	// protects pools
	mutex sync.Mutex
}

// NewPools creates ConnPools
func NewPools() *ConnPools {
	c := &ConnPools{}
	c.pools.Store(make(map[string]ConnPool))
	return c
}

// Get gets ConnPool by name
func (c *ConnPools) Get(name string) (ConnPool, bool) {
	pool, ok := c.pools.Load().(map[string]ConnPool)[name]
	return pool, ok
}

// Set stores ConnPool
func (c *ConnPools) Set(connPool ConnPool) {
	c.mutex.Lock()
	pools := c.pools.Load().(map[string]ConnPool)
	m := make(map[string]ConnPool, len(pools)+1)
	for k, v := range pools {
		m[k] = v
	}
	m[connPool.Name()] = connPool
	c.pools.Store(m)
	c.mutex.Unlock()
}

// Del delects ConnPool by name, and close the ConnPool.
func (c *ConnPools) Del(name string) {
	c.mutex.Lock()
	pools := c.pools.Load().(map[string]ConnPool)
	m := make(map[string]ConnPool, len(pools))
	for k, v := range pools {
		if k == name {
			v.Close()
		} else {
			m[k] = v
		}
	}
	c.pools.Store(m)
	c.mutex.Unlock()
}
