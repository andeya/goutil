# goutil [![report card](https://goreportcard.com/badge/github.com/henrylee2cn/goutil?style=flat-square)](http://goreportcard.com/report/henrylee2cn/goutil) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/henrylee2cn/goutil)

Common and useful utils for the Go project development.

## 1. Inclusion criteria

- Only rely on the Go standard package
- Functions or lightweight packages
- Non-business related general tools

## 2. Contents

- [Calendar](#calendar) Chinese Lunar Calendar, Solar Calendar and cron time rules
- [CoarseTime](#coarsetime) Current time truncated to the nearest 100ms
- [Errors](#errors) Improved errors package.
- [Graceful](#graceful) Shutdown or reboot current process gracefully.
- [GoPool](#gopool) Goroutines' pool
- [ResPool](#respool) Resources' pool
- [Workshop](#workshop) working workshop
- [Various](#various) Various small functions


## 3. UtilsAPI

### Calendar

Chinese Lunar Calendar, Solar Calendar and cron time rules.

- import it

	```go
	"github.com/henrylee2cn/goutil/calendar"
	```

[Calendar details](calendar/README.md)

### CoarseTime

The current time truncated to the nearest second.

- import it

	```go
	"github.com/henrylee2cn/goutil/coarsetime"
	```

- FloorTimeNow returns the current time from the range (now-100ms,now].
This is a faster alternative to time.Now().

	```go
	func FloorTimeNow() time.Time
	```

- CeilingTimeNow returns the current time from the range [now,now+100ms).
This is a faster alternative to time.Now().

	```go
	func CeilingTimeNow() time.Time
	```

### Errors

Errors is improved errors package.

- import it

	```go
	"github.com/henrylee2cn/goutil/errors"
	```

- New returns an error that formats as the given text.

	```go
	func New(text string) error
	```

- Errorf formats according to a format specifier and returns the string as a value that satisfies error.

	```go
	func Errorf(format string, a ...interface{}) error
	```

- Merge merges multi errors.

	```go
	func Merge(errs ...error) error
	```

- Append appends multiple errors to the error.

	```go
	func Append(err error, errs ...error) error
	```

### Graceful

Shutdown or reboot current process gracefully.

- import it

	```go
	"github.com/henrylee2cn/goutil/graceful"
	```

- GraceSignal open graceful shutdown or reboot signal.

	```go
	func GraceSignal()
	```

- SetShutdown sets the function which is called after the process shutdown,
and the time-out period for the process shutdown.
If 0<=timeout<5s, automatically use 'MinShutdownTimeout'(5s).
If timeout<0, indefinite period.
'firstSweepFunc' is first executed.
'beforeExitingFunc' is executed before process exiting.

	```go
	func SetShutdown(timeout time.Duration, firstSweepFunc, beforeExitingFunc func() error)
	```

- Shutdown closes all the frame process gracefully.
Parameter timeout is used to reset time-out period for the process shutdown.

	```go
	func Shutdown(timeout ...time.Duration)
	```

- Reboot all the frame process gracefully.
Notes: Windows system are not supported!

	```go
	func Reboot(timeout ...time.Duration)
	```

- AddInherited adds the files and envs to be inherited by the new process.
Notes:
 Only for reboot!
 Windows system are not supported!

	```go
	func AddInherited(procFiles []*os.File, envs []*Env)
	```

- Logger logger interface

	```go
	Logger interface {
		Infof(format string, v ...interface{})
		Errorf(format string, v ...interface{})
	}
	```

- SetLog resets logger

	```go
	func SetLog(logger Logger)
	```

### GoPool

GoPool is a Goroutines pool. It can control concurrent numbers, reuse goroutines.

- import it

	```go
	"github.com/henrylee2cn/goutil/pool"
	```

- GoPool executes concurrently incoming function via a pool of goroutines in
FILO order, i.e. the most recently stopped goroutine will execute the next
incoming function.
Such a scheme keeps CPU caches hot (in theory).

	```go
	type GoPool struct {
		// Has unexported fields.
	}
	```
	    
- NewGoPool creates a new *GoPool.
If maxGoroutinesAmount<=0, will use default value.
If maxGoroutineIdleDuration<=0, will use default value.

	```go
	func NewGoPool(maxGoroutinesAmount int, maxGoroutineIdleDuration time.Duration) *GoPool
	```

- Go executes the function via a goroutine.
If returns non-nil, the function cannot be executed because exceeded maxGoroutinesAmount limit.

	```go
	func (gp *GoPool) Go(fn func()) error
	```

- TryGo tries to execute the function via goroutine.
If there are no concurrent resources, execute it synchronously.

	```go
	func (gp *GoPool) TryGo(fn func())
	```

- Stop starts GoPool.
If calling 'Go' after calling 'Stop', will no longer reuse goroutine.

	```go
	func (gp *GoPool) Stop()
	```

### ResPool

ResPool is a high availability/high concurrent resource pool, which automatically manages the number of resources.
So it is similar to database/sql's db pool.

- import it

	```go
	"github.com/henrylee2cn/goutil/pool"
	```

- ResPool is a pool of zero or more underlying avatar(resource).
It's safe for concurrent use by multiple goroutines.
ResPool creates and frees resource automatically;
it also maintains a free pool of idle avatar(resource).

	```go
	type ResPool interface {
		// Name returns the name.
		Name() string
		// Get returns a object in Resource type.
		Get() (Resource, error)
		// GetContext returns a object in Resource type.
		// Support context cancellation.
		GetContext(context.Context) (Resource, error)
		// Put gives a resource back to the ResPool.
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
		// Close closes the ResPool, releasing any open resources.
		//
		// It is rare to close a ResPool, as the ResPool handle is meant to be
		// long-lived and shared between many goroutines.
		Close() error
		// Stats returns resource statistics.
		Stats() ResPoolStats
	}
	```

- NewResPool creates ResPool.

	```go
	func NewResPool(name string, newfunc func(context.Context) (Resource, error)) ResPool
	```

- Resource is a resource that can be stored in the ResPool.

	```go
	type Resource interface {
		// SetAvatar stores the contact with resPool
		// Do not call it yourself, it is only called by (*ResPool).get, and will only be called once
		SetAvatar(*Avatar)
		// GetAvatar gets the contact with resPool
		// Do not call it yourself, it is only called by (*ResPool).Put
		GetAvatar() *Avatar
		// Close closes the original source
		// No need to call it yourself, it is only called by (*Avatar).close
		Close() error
	}
	```

- Avatar links a Resource with a mutex, to be held during all calls into the Avatar.

	```go
	type Avatar struct {
		// Has unexported fields.
	}
	```

- Free releases self to the ResPool.
If error is not nil, close it.

	```go
	func (avatar *Avatar) Free(err error)
	```

- ResPool returns ResPool to which it belongs.

	```go
	func (avatar *Avatar) ResPool() ResPool
	```

- ResPools stores ResPool.

	```go
	type ResPools struct {
		// Has unexported fields.
	}
	```

- NewResPools creates a new ResPools.

	```go
	func NewResPools() *ResPools
	```

- Clean delects and close all the ResPools.

	```go
	func (c *ResPools) Clean()
	```

- Del delects ResPool by name, and close the ResPool.

	```go
	func (c *ResPools) Del(name string)
	```

- Get gets ResPool by name.

	```go
	func (c *ResPools) Get(name string) (ResPool, bool)
	```

- GetAll gets all the ResPools.

	```go
	func (c *ResPools) GetAll() []ResPool
	```

- Set stores ResPool.
If the same name exists, will close and cover it.

	```go
	func (c *ResPools) Set(pool ResPool)
	```

### Workshop

- import it

	```go
	"github.com/henrylee2cn/goutil/pool"
	```

- Type definition

	```go
	type (
		// Worker woker interface
		// Note: Worker can not be implemented using empty structures(struct{})!
		Worker interface {
			Health() bool
			Close() error
		}
		// Workshop working workshop
		Workshop struct {
			// Has unexported fields.
		}
	)
	```
	
- NewWorkshop creates a new workshop.
<br>If maxQuota<=0, will use default value.
<br>If maxIdleDuration<=0, will use default value.
<br>Note: Worker can not be implemented using empty structures(struct{})!

	```go
	func NewWorkshop(maxQuota int, maxIdleDuration time.Duration, newWorkerFunc func() (Worker, error)) *Workshop
	```

- Close wait for all the work to be completed and close the workshop.

	```go
	func (w *Workshop) Close()
	```

- Callback assigns a healthy worker to execute the function.

	```go
	func (w *Workshop) Callback(fn func(Worker) error) error
	```

- Fire marks the worker to reduce a job.
<br>If the worker does not belong to the workshop, close the worker.

	```go
	func (w *Workshop) Fire(worker Worker)
	```

- Hire hires a healthy worker and marks the worker to add a job.

	```go
	func (w *Workshop) Hire() (Worker, error)
	```

- Stats returns the current workshop stats.

	```go
	func (w *Workshop) Stats() *WorkshopStats
	```

### Various

Various small functions.

- import it

	```go
	"github.com/henrylee2cn/goutil"
	```

- BytesToString convert []byte type to string type.

	```go
	func BytesToString(b []byte) string
	```

- StringToBytes convert string type to []byte type.
NOTE: panic if modify the member value of the []byte.

	```go
	func StringToBytes(s string) []byte
	```

- RandomBytes returns securely generated random bytes. It will panic if the system's secure random number generator fails to function correctly.

	```go
	func RandomBytes(n int) []byte
	```

- RandomString returns a URL-safe, base64 encoded securely generated random string. It will panic if the system's secure random number generator fails to function correctly.
The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.

	```go
	func RandomString(n int) string
	```

- CamelString converts the accepted string to a camel string (xx_yy to XxYy)

	```go
	func CamelString(s string) string
	```

- SnakeString converts the accepted string to a snake string (XxYy to xx_yy)

	```go
	func SnakeString(s string) string
	```

- ObjectName gets the type name of the object

	```go
	func ObjectName(obj interface{}) string
	```

- JsQueryEscape escapes the string in javascript standard so it can be safely placed inside a URL query.

	```go
	func JsQueryEscape(s string) string
	```

- JsQueryUnescape does the inverse transformation of JsQueryEscape, converting %AB into the byte 0xAB and '+' into ' ' (space). It returns an error if any % is not followed by two hexadecimal digits.

	```go
	func JsQueryUnescape(s string) (string, error)
	```

- Map is a concurrent map with loads, stores, and deletes.
It is safe for multiple goroutines to call a Map's methods concurrently.

	```go
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
		// Range calls f sequentially for each key and value present in the map.
		// If f returns false, range stops the iteration.
		Range(f func(key, value interface{}) bool)
		// Random returns a pair kv randomly.
		// If exist=false, no kv data is exist.
		Random() (key, value interface{}, exist bool)
		// Delete deletes the value for a key.
		Delete(key interface{})
		// Clear clears all current data in the map.
		Clear()
		// Len returns the length of the map.
		Len() int
	}
	```

- RwMap creates a new concurrent safe map with sync.RWMutex.
The normal Map is high-performance mapping under low concurrency conditions.

	```go
	func RwMap(capacity ...int) Map
	```

- AtomicMap creates a concurrent map with amortized-constant-time loads, stores, and deletes.
It is safe for multiple goroutines to call a atomicMap's methods concurrently.
From go v1.9 sync.Map.

	```go
	func AtomicMap() Map
	```

- SelfPath gets compiled executable file absolute path.

	```go
	func SelfPath() string
	```

- SelfDir gets compiled executable file directory.

	```go
	func SelfDir() string
	```

- RelPath gets relative path.

	```go
	func RelPath(targpath string) string
	```

- SelfChdir switch the working path to my own path.

	```go
	func SelfChdir()
	```

- FileExists reports whether the named file or directory exists.

	```go
	func FileExists(name string) bool
	```

- SearchFile Search a file in paths.
This is often used in search config file in `/etc` `~/`

	```go
	func SearchFile(filename string, paths ...string) (fullpath string, err error)
	```

- GrepFile like command grep -E.
For example: GrepFile(`^hello`, "hello.txt").
`\n` is striped while read

	```go
	func GrepFile(patten string, filename string) (lines []string, err error)
	```

- WalkDirs traverses the directory, return to the relative path.
You can specify the suffix.

	```go
	func WalkDirs(targpath string, suffixes ...string) (dirlist []string)
	```

- IsExportedOrBuiltinType is this type exported or a builtin?

	```go
	func IsExportedOrBuiltinType(t reflect.Type) bool
	```

- IsExportedName is this an exported - upper case - name?

	```go
	func IsExportedName(name string) bool
	```

- PanicTrace trace panic stack info.

	```go
	func PanicTrace(kb int) []byte
	```

- ExtranetIP get external IP addr.

	```go
	func ExtranetIP() (ip string, err error)
	```

- IntranetIP get internal IP addr.

	```go
	func IntranetIP() (string, error)
	```
