# tpack

Go underlying type data.

- import it

    ```go
    "github.com/henrylee2cn/goutil/tpack"
    ```

- doc

    ```go
    // T go underlying type data
    type T struct {
        // Has unexported fields.
    }

    // Unpack unpack i to go underlying type data.
    func Unpack(i interface{}) T

    // RuntimeTypeID gets the underlying type ID in current runtime.
    // NOTE:
    //  *A and A gets the same runtime type ID;
    //  It is 10 times performance of reflect.TypeOf(i).String().
    func (t T) RuntimeTypeID() int32

    // Kind gets the reflect.Kind fastly.
    func (t T) Kind() reflect.Kind

    // Pointer gets the pointer of i.
    // NOTE:
    //  *A and A, gets diffrent pointer
    func (t T) Pointer() uintptr

    // TypeOf is equivalent to reflect.TypeOf.
    func (t T) TypeOf() reflect.Type

    // ValueOf is equivalent to reflect.ValueOf.
    func (t T) ValueOf() reflect.Value

    // RuntimeTypeID gets the underlying type ID in current runtime from reflect.Type.
    // NOTE:
    //  *A and A gets the same runtime type ID;
    //  It is 10 times performance of t.String().
    func RuntimeTypeID(t reflect.Type) int32
    ```

## Benchmark

```
goos: darwin
goarch: amd64
pkg: github.com/henrylee2cn/goutil/tpack
BenchmarkUnpack_tpack-4   	2000000000	         0.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkTypeOf_go-4      	200000000	        10.3 ns/op	       0 B/op	       0 allocs/op
```