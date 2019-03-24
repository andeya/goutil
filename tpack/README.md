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

    // TypeID returns the underlying type ID.
    // It is 60 times performance of reflect.TypeOf(i).String()
    func (t T) TypeID() int32

    // TypeOf is equivalent to reflect.TypeOf.
    func (t T) TypeOf() reflect.Type

    // ValueOf is equivalent to reflect.ValueOf.
    func (t T) ValueOf() reflect.Value

    // TypeID get underlying type ID from reflect.Type.
    // It is 60 times performance of t.String()
    func TypeID(t reflect.Type) int32
    ```

## Benchmark

```
goos: darwin
goarch: amd64
pkg: github.com/henrylee2cn/goutil/tpack
BenchmarkUnpack_tpack-4   	2000000000	         0.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkValueOf_go-4     	30000000	        51.0 ns/op	      16 B/op	       1 allocs/op
```