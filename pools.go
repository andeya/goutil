package srcpool

import (
	"sort"
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

// GetAll gets all the Pools
func (c *Pools) GetAll() []Pool {
	all := c.pools.Load().(map[string]Pool)
	pools := make(pools, 0, len(all))
	for _, pool := range all {
		pools = append(pools, pool)
	}
	sort.Sort(pools)
	return pools
}

// Set stores Pool.
// If the same name exists, will close and cover it.
func (c *Pools) Set(pool Pool) {
	c.mutex.Lock()
	pools := c.pools.Load().(map[string]Pool)
	m := make(map[string]Pool, len(pools)+1)
	name := pool.Name()
	for k, v := range pools {
		if k == name {
			v.Close()
		} else {
			m[k] = v
		}
	}
	m[name] = pool
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

// Clean delects and close all the Pools.
func (c *Pools) Clean() {
	c.mutex.Lock()
	pools := c.pools.Load().(map[string]Pool)
	for _, v := range pools {
		v.Close()
	}
	c.pools.Store(make(map[string]Pool))
	c.mutex.Unlock()
}

type pools []Pool

func (p pools) Len() int {
	return len(p)
}

func (p pools) Less(i, j int) bool {
	return p[i].Name() < p[j].Name()
}

func (p pools) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
