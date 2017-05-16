# srcpool   [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/srcpool)

srcpool is a high availability / high concurrent resource pool.
It automatically manages the number of resources, which is similar to database/sql's db pool.

```go
// NewFunc creates a new avatar.
//
// NewFunc may return a cached resource (one previously
// closed), but doing so is unnecessary; the pool package
// maintains a pool of idle resources for efficient re-use.
//
// The returned resource is only used by one goroutine at a
// time.
type NewFunc func(context.Context) (Resource, error)
```

```go
// Resource is a resource that can be stored in the Pool.
type Resource interface {
    // SetAvatar stores the contact with pool
    // Do not call it yourself, it is only called by (*Pool).get, and will only be called once
    SetAvatar(*Avatar)
    // GetAvatar gets the contact with pool
    // Do not call it yourself, it is only called by (*Pool).Put
    GetAvatar() *Avatar
    // Close closes the original source
    // No need to call it yourself, it is only called by (*Avatar).close
    Close() error
}
```

```go

// Pool is a pool of zero or more underlying avatar(resource).
// It's safe for concurrent use by multiple goroutines.
//
// Pool creates and frees resource automatically;
// it also maintains a free pool of idle avatar(resource).
type Pool interface {
    // Name returns the name.
    Name() string
    // Get returns a object in Resource type.
    Get() (Resource, error)
    // GetContext returns a object in Resource type.
    // Support context cancellation.
    GetContext(context.Context) (Resource, error)
    // Put gives a resource back to the Pool.
    // If error is not nil, close the avatar.
    Put(Resource, error)
    // Callback callbacks your handle function, returns the error of getting resource or handling.
    // Support recover panic.
    Callback(func(Resource) error) error
    // Callback callbacks your handle function, returns the error of getting resource or handling.
    // Support recover panic and context cancellation.
    CallbackContext(context.Context, func(Resource) error) error
    // SetMaxLifetime sets the maximum amount of time a resource may be reused.
    //
    // Expired resource may be closed lazily before reuse.
    //
    // If d <= 0, resource are reused forever.
    SetMaxLifetime(d time.Duration)
    // SetMaxIdle sets the maximum number of resources in the idle
    // resource pool.
    //
    // If SetMaxIdle is greater than 0 but less than the new MaxIdle
    // then the new MaxIdle will be reduced to match the SetMaxIdle limit
    //
    // If n <= 0, no idle resources are retained.
    SetMaxIdle(n int)
    // SetMaxOpen sets the maximum number of open resources.
    //
    // If MaxIdle is greater than 0 and the new MaxOpen is less than
    // MaxIdle, then MaxIdle will be reduced to match the new
    // MaxOpen limit
    //
    // If n <= 0, then there is no limit on the number of open resources.
    // The default is 0 (unlimited).
    SetMaxOpen(n int)
    // Close closes the Pool, releasing any open resources.
    //
    // It is rare to close a Pool, as the Pool handle is meant to be
    // long-lived and shared between many goroutines.
    Close() error
    // Stats returns resource statistics.
    Stats() PoolStats
}
```