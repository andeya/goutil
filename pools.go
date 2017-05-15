package srcpool

import (
	"sync"
	"sync/atomic"
)

// Pools stores Pool
type Pools struct {
	// stores 'map[string]Pool',
	// one server node has one connection pool.
	pools atomic.Value
	// protects pools
	mutex sync.Mutex
}

// NewPools creates Pools
func NewPools() *Pools {
	c := &Pools{}
	c.pools.Store(make(map[string]Pool))
	return c
}

// Get gets Pool by name
func (c *Pools) Get(name string) (Pool, bool) {
	pool, ok := c.pools.Load().(map[string]Pool)[name]
	return pool, ok
}

// Set stores Pool.
// If the same name exists, will close and cover it.
func (c *Pools) Set(connPool Pool) {
	c.mutex.Lock()
	pools := c.pools.Load().(map[string]Pool)
	m := make(map[string]Pool, len(pools)+1)
	name := connPool.Name()
	for k, v := range pools {
		if k == name {
			v.Close()
		} else {
			m[k] = v
		}
	}
	m[name] = connPool
	c.pools.Store(m)
	c.mutex.Unlock()
}

// Del delects Pool by name, and close the Pool.
func (c *Pools) Del(name string) {
	c.mutex.Lock()
	pools := c.pools.Load().(map[string]Pool)
	m := make(map[string]Pool, len(pools))
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
