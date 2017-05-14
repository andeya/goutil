package connpool

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ConnPool is a pool of zero or more underlying connections.
// It's safe for concurrent use by multiple goroutines.
//
// ConnPool creates and frees connections automatically;
// it also maintains a free pool of idle connections.
type ConnPool interface {
	// Name returns the connPool's name.
	Name() string
	// New creates a new net.Conn.
	New() (net.Conn, error)
	// Get returns a net.Conn.
	Get() (net.Conn, error)
	// GetContext returns a net.Conn, support context cancellation.
	GetContext(ctx context.Context) (net.Conn, error)
	// Put adds a net.Conn to the connPool's free pool.
	Put(conn net.Conn)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are reused forever.
	SetConnMaxLifetime(d time.Duration)
	// SetMaxIdleConns sets the maximum number of connections in the idle
	// connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
	//
	// If n <= 0, no idle connections are retained.
	SetMaxIdleConns(n int)
	// SetMaxOpenConns sets the maximum number of open connections to the net.
	//
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	SetMaxOpenConns(n int)
	// Close closes the ConnPool, releasing any open resources.
	//
	// It is rare to Close a ConnPool, as the ConnPool handle is meant to be
	// long-lived and shared between many goroutines.
	Close() error
	// Stats returns connection statistics.
	Stats() ConnPoolStats
}

// This is the size of the connectionOpener request chan (ConnPool.openerCh).
// This value should be larger than the maximum typical value
// used for connPool.maxOpen. If maxOpen is significantly larger than
// connectionRequestQueueSize then it is possible for ALL calls into the *ConnPool
// to block until the connectionOpener can satisfy the backlog of requests.
const connectionRequestQueueSize = 1000000

// New creates ConnPool.
func New(name string, dialFunc DialFunc) (ConnPool, error) {
	connPool := &pool{
		driver:       dialFunc,
		name:         name,
		openerCh:     make(chan struct{}, connectionRequestQueueSize),
		lastPut:      make(map[*driverConn]string),
		connRequests: make(map[uint64]chan connRequest),
	}
	go connPool.connectionOpener()
	return connPool, nil
}

// DialFunc returns a new connection.
//
// Dial may return a cached connection (one previously
// closed), but doing so is unnecessary; the connpool package
// maintains a pool of idle connections for efficient re-use.
//
// The returned connection is only used by one goroutine at a
// time.
type DialFunc func() (net.Conn, error)

type pool struct {
	driver DialFunc
	name   string
	// numClosed is an atomic counter which represents a total number of
	// closed connections. Stmt.openStmt checks it before cleaning closed
	// connections in Stmt.css.
	numClosed uint64

	mu           sync.Mutex // protects following fields
	freeConn     []*driverConn
	connRequests map[uint64]chan connRequest
	nextRequest  uint64 // Next key to use in connRequests.
	numOpen      int    // number of opened and pending open connections
	// Used to signal the need for new connections
	// a goroutine running connectionOpener() reads on this chan and
	// maybeOpenNewConnections sends on the chan (one send per needed connection)
	// It is closed during connPool.Close(). The close tells the connectionOpener
	// goroutine to exit.
	openerCh    chan struct{}
	closed      bool
	dep         map[finalCloser]depSet
	lastPut     map[*driverConn]string // stacktrace of last conn's put; debug only
	maxIdle     int                    // zero means defaultMaxIdleConns; negative means 0
	maxOpen     int                    // <= 0 means unlimited
	maxLifetime time.Duration          // maximum amount of time a connection may be reused
	cleanerCh   chan struct{}
}

// connReuseStrategy determines how (*ConnPool).conn returns net connections.
type connReuseStrategy uint8

const (
	// alwaysNewConn forces a new connection to the net.
	alwaysNewConn connReuseStrategy = iota
	// cachedOrNewConn returns a cached connection, if available, else waits
	// for one to become available (if MaxOpenConns has been reached) or
	// creates a new net connection.
	cachedOrNewConn
)

// driverConn wraps a driver.Conn with a mutex, to
// be held during all calls into the Conn. (including any calls onto
// interfaces returned via that Conn, such as calls on Tx, Stmt,
// Result, Rows)
type driverConn struct {
	connPool  *pool
	createdAt time.Time

	sync.Mutex // guards following
	net.Conn
	closed      bool
	finalClosed bool // Conn.Close has been called

	// guarded by connPool.mu
	inUse bool
}

func (dc *driverConn) releaseConn(err error) {
	dc.connPool.putConn(dc, err)
}

func (dc *driverConn) expired(timeout time.Duration) bool {
	if timeout <= 0 {
		return false
	}
	return dc.createdAt.Add(timeout).Before(nowFunc())
}

// the dc.connPool's Mutex is held.
func (dc *driverConn) closeConnPoolLocked() func() error {
	dc.Lock()
	defer dc.Unlock()
	if dc.closed {
		return func() error { return errors.New("connpool: duplicate driverConn close") }
	}
	dc.closed = true
	return dc.connPool.removeDepLocked(dc, dc)
}

func (dc *driverConn) Close() error {
	dc.Lock()
	if dc.closed {
		dc.Unlock()
		return errors.New("connpool: duplicate driverConn close")
	}
	dc.closed = true
	dc.Unlock() // not defer; removeDep finalClose calls may need to lock

	// And now updates that require holding dc.mu.Lock.
	dc.connPool.mu.Lock()
	fn := dc.connPool.removeDepLocked(dc, dc)
	dc.connPool.mu.Unlock()
	return fn()
}

func (dc *driverConn) finalClose() error {
	var err error
	withLock(dc, func() {
		dc.finalClosed = true
		err = dc.Conn.Close()
		dc.Conn = nil
	})

	dc.connPool.mu.Lock()
	dc.connPool.numOpen--
	dc.connPool.maybeOpenNewConnections()
	dc.connPool.mu.Unlock()

	atomic.AddUint64(&dc.connPool.numClosed, 1)
	return err
}

// depSet is a finalCloser's outstanding dependencies
type depSet map[interface{}]bool // set of true bools

// The finalCloser interface is used by (*ConnPool).addDep and related
// dependency reference counting.
type finalCloser interface {
	// finalClose is called when the reference count of an object
	// goes to zero. (*ConnPool).mu is not held while calling it.
	finalClose() error
}

// addDep notes that x now depends on dep, and x's finalClose won't be
// called until all of x's dependencies are removed with removeDep.
func (connPool *pool) addDep(x finalCloser, dep interface{}) {
	//println(fmt.Sprintf("addDep(%T %p, %T %p)", x, x, dep, dep))
	connPool.mu.Lock()
	defer connPool.mu.Unlock()
	connPool.addDepLocked(x, dep)
}

func (connPool *pool) addDepLocked(x finalCloser, dep interface{}) {
	if connPool.dep == nil {
		connPool.dep = make(map[finalCloser]depSet)
	}
	xdep := connPool.dep[x]
	if xdep == nil {
		xdep = make(depSet)
		connPool.dep[x] = xdep
	}
	xdep[dep] = true
}

// removeDep notes that x no longer depends on dep.
// If x still has dependencies, nil is returned.
// If x no longer has any dependencies, its finalClose method will be
// called and its error value will be returned.
func (connPool *pool) removeDep(x finalCloser, dep interface{}) error {
	connPool.mu.Lock()
	fn := connPool.removeDepLocked(x, dep)
	connPool.mu.Unlock()
	return fn()
}

func (connPool *pool) removeDepLocked(x finalCloser, dep interface{}) func() error {
	//println(fmt.Sprintf("removeDep(%T %p, %T %p)", x, x, dep, dep))

	xdep, ok := connPool.dep[x]
	if !ok {
		panic(fmt.Sprintf("unpaired removeDep: no deps for %T", x))
	}

	l0 := len(xdep)
	delete(xdep, dep)

	switch len(xdep) {
	case l0:
		// Nothing removed. Shouldn't happen.
		panic(fmt.Sprintf("unpaired removeDep: no %T dep on %T", dep, x))
	case 0:
		// No more dependencies.
		delete(connPool.dep, x)
		return x.finalClose
	default:
		// Dependencies remain.
		return func() error { return nil }
	}
}

// Close closes the ConnPool, releasing any open resources.
//
// It is rare to Close a ConnPool, as the ConnPool handle is meant to be
// long-lived and shared between many goroutines.
func (connPool *pool) Close() error {
	connPool.mu.Lock()
	if connPool.closed { // Make ConnPool.Close idempotent
		connPool.mu.Unlock()
		return nil
	}
	close(connPool.openerCh)
	if connPool.cleanerCh != nil {
		close(connPool.cleanerCh)
	}
	var err error
	fns := make([]func() error, 0, len(connPool.freeConn))
	for _, dc := range connPool.freeConn {
		fns = append(fns, dc.closeConnPoolLocked())
	}
	connPool.freeConn = nil
	connPool.closed = true
	for _, req := range connPool.connRequests {
		close(req)
	}
	connPool.mu.Unlock()
	for _, fn := range fns {
		err1 := fn()
		if err1 != nil {
			err = err1
		}
	}
	return err
}

const defaultMaxIdleConns = 2

func (connPool *pool) maxIdleConnsLocked() int {
	n := connPool.maxIdle
	switch {
	case n == 0:
		// TODO(bradfitz): ask driver, if supported, for its default preference
		return defaultMaxIdleConns
	case n < 0:
		return 0
	default:
		return n
	}
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
//
// If n <= 0, no idle connections are retained.
func (connPool *pool) SetMaxIdleConns(n int) {
	connPool.mu.Lock()
	if n > 0 {
		connPool.maxIdle = n
	} else {
		// No idle connections.
		connPool.maxIdle = -1
	}
	// Make sure maxIdle doesn't exceed maxOpen
	if connPool.maxOpen > 0 && connPool.maxIdleConnsLocked() > connPool.maxOpen {
		connPool.maxIdle = connPool.maxOpen
	}
	var closing []*driverConn
	idleCount := len(connPool.freeConn)
	maxIdle := connPool.maxIdleConnsLocked()
	if idleCount > maxIdle {
		closing = connPool.freeConn[maxIdle:]
		connPool.freeConn = connPool.freeConn[:maxIdle]
	}
	connPool.mu.Unlock()
	for _, c := range closing {
		c.Close()
	}
}

// SetMaxOpenConns sets the maximum number of open connections to the net.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func (connPool *pool) SetMaxOpenConns(n int) {
	connPool.mu.Lock()
	connPool.maxOpen = n
	if n < 0 {
		connPool.maxOpen = 0
	}
	syncMaxIdle := connPool.maxOpen > 0 && connPool.maxIdleConnsLocked() > connPool.maxOpen
	connPool.mu.Unlock()
	if syncMaxIdle {
		connPool.SetMaxIdleConns(n)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are reused forever.
func (connPool *pool) SetConnMaxLifetime(d time.Duration) {
	if d < 0 {
		d = 0
	}
	connPool.mu.Lock()
	// wake cleaner up when lifetime is shortened.
	if d > 0 && d < connPool.maxLifetime && connPool.cleanerCh != nil {
		select {
		case connPool.cleanerCh <- struct{}{}:
		default:
		}
	}
	connPool.maxLifetime = d
	connPool.startCleanerLocked()
	connPool.mu.Unlock()
}

// Name returns the connPool's name.
func (connPool *pool) Name() string {
	return connPool.name
}

// New creates a new net.Conn.
func (connPool *pool) New() (net.Conn, error) {
	return connPool.driver()
}

// startCleanerLocked starts connectionCleaner if needed.
func (connPool *pool) startCleanerLocked() {
	if connPool.maxLifetime > 0 && connPool.numOpen > 0 && connPool.cleanerCh == nil {
		connPool.cleanerCh = make(chan struct{}, 1)
		go connPool.connectionCleaner(connPool.maxLifetime)
	}
}

func (connPool *pool) connectionCleaner(d time.Duration) {
	const minInterval = time.Second

	if d < minInterval {
		d = minInterval
	}
	t := time.NewTimer(d)

	for {
		select {
		case <-t.C:
		case <-connPool.cleanerCh: // maxLifetime was changed or connPool was closed.
		}

		connPool.mu.Lock()
		d = connPool.maxLifetime
		if connPool.closed || connPool.numOpen == 0 || d <= 0 {
			connPool.cleanerCh = nil
			connPool.mu.Unlock()
			return
		}

		expiredSince := nowFunc().Add(-d)
		var closing []*driverConn
		for i := 0; i < len(connPool.freeConn); i++ {
			c := connPool.freeConn[i]
			if c.createdAt.Before(expiredSince) {
				closing = append(closing, c)
				last := len(connPool.freeConn) - 1
				connPool.freeConn[i] = connPool.freeConn[last]
				connPool.freeConn[last] = nil
				connPool.freeConn = connPool.freeConn[:last]
				i--
			}
		}
		connPool.mu.Unlock()

		for _, c := range closing {
			c.Close()
		}

		if d < minInterval {
			d = minInterval
		}
		t.Reset(d)
	}
}

// ConnPoolStats contains net statistics.
type ConnPoolStats struct {
	// OpenConnections is the number of open connections to the net.
	OpenConnections int
}

// Stats returns connection statistics.
func (connPool *pool) Stats() ConnPoolStats {
	connPool.mu.Lock()
	stats := ConnPoolStats{
		OpenConnections: connPool.numOpen,
	}
	connPool.mu.Unlock()
	return stats
}

// Assumes connPool.mu is locked.
// If there are connRequests and the connection limit hasn't been reached,
// then tell the connectionOpener to open new connections.
func (connPool *pool) maybeOpenNewConnections() {
	numRequests := len(connPool.connRequests)
	if connPool.maxOpen > 0 {
		numCanOpen := connPool.maxOpen - connPool.numOpen
		if numRequests > numCanOpen {
			numRequests = numCanOpen
		}
	}
	for numRequests > 0 {
		connPool.numOpen++ // optimistically
		numRequests--
		if connPool.closed {
			return
		}
		connPool.openerCh <- struct{}{}
	}
}

// Runs in a separate goroutine, opens new connections when requested.
func (connPool *pool) connectionOpener() {
	for range connPool.openerCh {
		connPool.openNewConnection()
	}
}

// Open one new connection
func (connPool *pool) openNewConnection() {
	// maybeOpenNewConnctions has already executed connPool.numOpen++ before it sent
	// on connPool.openerCh. This function must execute connPool.numOpen-- if the
	// connection fails or is closed before returning.
	ci, err := connPool.driver()
	connPool.mu.Lock()
	defer connPool.mu.Unlock()
	if connPool.closed {
		if err == nil {
			ci.Close()
		}
		connPool.numOpen--
		return
	}
	if err != nil {
		connPool.numOpen--
		connPool.putConnConnPoolLocked(nil, err)
		connPool.maybeOpenNewConnections()
		return
	}
	dc := &driverConn{
		connPool:  connPool,
		createdAt: nowFunc(),
		Conn:      ci,
	}
	if connPool.putConnConnPoolLocked(dc, err) {
		connPool.addDepLocked(dc, dc)
	} else {
		connPool.numOpen--
		ci.Close()
	}
}

// connRequest represents one request for a new connection
// When there are no idle connections available, ConnPool.conn will create
// a new connRequest and put it on the connPool.connRequests list.
type connRequest struct {
	conn *driverConn
	err  error
}

var errConnPoolClosed = errors.New("connpool: net is closed")

// nextRequestKeyLocked returns the next connection request key.
// It is assumed that nextRequest will not overflow.
func (connPool *pool) nextRequestKeyLocked() uint64 {
	next := connPool.nextRequest
	connPool.nextRequest++
	return next
}

var ErrBadConn = errors.New("connpool: bad connection")

// conn returns a newly-opened or cached *driverConn.
func (connPool *pool) conn(ctx context.Context, strategy connReuseStrategy) (*driverConn, error) {
	connPool.mu.Lock()
	if connPool.closed {
		connPool.mu.Unlock()
		return nil, errConnPoolClosed
	}
	// Check if the context is expired.
	select {
	default:
	case <-ctx.Done():
		connPool.mu.Unlock()
		return nil, ctx.Err()
	}
	lifetime := connPool.maxLifetime

	// Prefer a free connection, if possible.
	numFree := len(connPool.freeConn)
	if strategy == cachedOrNewConn && numFree > 0 {
		conn := connPool.freeConn[0]
		copy(connPool.freeConn, connPool.freeConn[1:])
		connPool.freeConn = connPool.freeConn[:numFree-1]
		conn.inUse = true
		connPool.mu.Unlock()
		if conn.expired(lifetime) {
			conn.Close()
			return nil, ErrBadConn
		}
		return conn, nil
	}

	// Out of free connections or we were asked not to use one. If we're not
	// allowed to open any more connections, make a request and wait.
	if connPool.maxOpen > 0 && connPool.numOpen >= connPool.maxOpen {
		// Make the connRequest channel. It's buffered so that the
		// connectionOpener doesn't block while waiting for the req to be read.
		req := make(chan connRequest, 1)
		reqKey := connPool.nextRequestKeyLocked()
		connPool.connRequests[reqKey] = req
		connPool.mu.Unlock()

		// Timeout the connection request with the context.
		select {
		case <-ctx.Done():
			// Remove the connection request and ensure no value has been sent
			// on it after removing.
			connPool.mu.Lock()
			delete(connPool.connRequests, reqKey)
			connPool.mu.Unlock()
			select {
			default:
			case ret, ok := <-req:
				if ok {
					connPool.putConn(ret.conn, ret.err)
				}
			}
			return nil, ctx.Err()
		case ret, ok := <-req:
			if !ok {
				return nil, errConnPoolClosed
			}
			if ret.err == nil && ret.conn.expired(lifetime) {
				ret.conn.Close()
				return nil, ErrBadConn
			}
			return ret.conn, ret.err
		}
	}

	connPool.numOpen++ // optimistically
	connPool.mu.Unlock()
	ci, err := connPool.driver()
	if err != nil {
		connPool.mu.Lock()
		connPool.numOpen-- // correct for earlier optimism
		connPool.maybeOpenNewConnections()
		connPool.mu.Unlock()
		return nil, err
	}
	connPool.mu.Lock()
	dc := &driverConn{
		connPool:  connPool,
		createdAt: nowFunc(),
		Conn:      ci,
	}
	connPool.addDepLocked(dc, dc)
	dc.inUse = true
	connPool.mu.Unlock()
	return dc, nil
}

// maxBadConnRetries is the number of maximum retries if the driver returns
// ErrBadConn to signal a broken connection before forcing a new
// connection to be opened.
const maxBadConnRetries = 2

// GetContext returns a net.Conn, support context cancellation.
func (connPool *pool) GetContext(ctx context.Context) (net.Conn, error) {
	var err error
	var ci *driverConn
	for i := 0; i < maxBadConnRetries; i++ {
		ci, err = connPool.conn(ctx, cachedOrNewConn)
		if err == nil {
			break
		}
	}
	if err != nil {
		return connPool.conn(ctx, alwaysNewConn)
	}
	return ci, err
}

// Get returns a net.Conn.
func (connPool *pool) Get() (net.Conn, error) {
	return connPool.GetContext(context.Background())
}

// Put adds a net.Conn to the connPool's free pool.
func (connPool *pool) Put(conn net.Conn) {
	if dc, ok := conn.(*driverConn); ok {
		connPool.putConn(dc, nil)
	}
}

// putConnHook is a hook for testing.
var putConnHook func(*pool, *driverConn)

// debugGetPut determines whether getConn & putConn calls' stack traces
// are returned for more verbose crashes.
const debugGetPut = false

// putConn adds a connection to the connPool's free pool.
// err is optionally the last error that occurred on this connection.
func (connPool *pool) putConn(dc *driverConn, err error) {
	connPool.mu.Lock()
	if !dc.inUse {
		if debugGetPut {
			fmt.Printf("putConn(%v) DUPLICATE was: %s\n\nPREVIOUS was: %s", dc, stack(), connPool.lastPut[dc])
		}
		panic("connpool: connection returned that was never out")
	}
	if debugGetPut {
		connPool.lastPut[dc] = stack()
	}
	dc.inUse = false

	if err != nil {
		// Don't reuse bad connections.
		// Since the conn is considered bad and is being discarded, treat it
		// as closed. Don't decrement the open count here, finalClose will
		// take care of that.
		connPool.maybeOpenNewConnections()
		connPool.mu.Unlock()
		dc.Close()
		return
	}
	if putConnHook != nil {
		putConnHook(connPool, dc)
	}
	added := connPool.putConnConnPoolLocked(dc, nil)
	connPool.mu.Unlock()

	if !added {
		dc.Close()
	}
}

// Satisfy a connRequest or put the driverConn in the idle pool and return true
// or return false.
// putConnConnPoolLocked will satisfy a connRequest if there is one, or it will
// return the *driverConn to the freeConn list if err == nil and the idle
// connection limit will not be exceeded.
// If err != nil, the value of dc is ignored.
// If err == nil, then dc must not equal nil.
// If a connRequest was fulfilled or the *driverConn was placed in the
// freeConn list, then true is returned, otherwise false is returned.
func (connPool *pool) putConnConnPoolLocked(dc *driverConn, err error) bool {
	if connPool.closed {
		return false
	}
	if connPool.maxOpen > 0 && connPool.numOpen > connPool.maxOpen {
		return false
	}
	if c := len(connPool.connRequests); c > 0 {
		var req chan connRequest
		var reqKey uint64
		for reqKey, req = range connPool.connRequests {
			break
		}
		delete(connPool.connRequests, reqKey) // Remove from pending requests.
		if err == nil {
			dc.inUse = true
		}
		req <- connRequest{
			conn: dc,
			err:  err,
		}
		return true
	} else if err == nil && !connPool.closed && connPool.maxIdleConnsLocked() > len(connPool.freeConn) {
		connPool.freeConn = append(connPool.freeConn, dc)
		connPool.startCleanerLocked()
		return true
	}
	return false
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], false)])
}

// withLock runs while holding lk.
func withLock(lk sync.Locker, fn func()) {
	lk.Lock()
	defer lk.Unlock() // in case fn panics
	fn()
}
