# tpack

Go underlying type data.

- import it

    ```go
    "github.com/henrylee2cn/goutil/tpack"
    ```

- doc

    ```go
    // U go underlying type data
    type U struct {
        // Has unexported fields.
    }

    // Unpack unpacks i to go underlying type data.
    func Unpack(i interface{}) U
   
    // From gets go underlying type data from reflect.Value.
    func From(v reflect.Value) U
    
    // RuntimeTypeID gets the underlying type ID in current runtime from reflect.Type.
    // NOTE:
    //  *A and A gets the same runtime type ID;
    //  It is 10 times performance of t.String().
    func RuntimeTypeID(t reflect.Type) int32
    
    // RuntimeTypeID gets the underlying type ID in current runtime.
    // NOTE:
    //  *A and A gets the same runtime type ID;
    //  It is 10 times performance of reflect.TypeOf(i).String().
    func (u U) RuntimeTypeID() int32

    // Kind gets the reflect.Kind fastly.
    func (u U) Kind() reflect.Kind

    // Elem returns the U that the interface i contains
    // or that the pointer i points to.
    func (u U) Elem() U

    // UnderlyingElem returns the underlying U that the interface i contains
    // or that the pointer i points to.
    func (u U) UnderlyingElem() U

    // Pointer gets the pointer of i.
    // NOTE:
    //  *T and T, gets diffrent pointer
    func (u U) Pointer() uintptr

    // IsNil reports whether its argument i is nil.
    func (u U) IsNil() bool

    // FuncForPC returns a *Func describing the function that contains the
    // given program counter address, or else nil.
    //
    // If pc represents multiple functions because of inlining, it returns
    // the a *Func describing the innermost function, but with an entry
    // of the outermost function.
    //
    // NOTE: Its kind must be a reflect.Func, otherwise it returns nil
    func (u U) FuncForPC() *runtime.Func
    ```

## Benchmark

```
goos: darwin
goarch: amd64
pkg: github.com/henrylee2cn/goutil/tpack
BenchmarkUnpack_tpack-4   	2000000000	         0.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkTypeOf_go-4      	200000000	        10.3 ns/op	       0 B/op	       0 allocs/op
```