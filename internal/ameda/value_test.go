package ameda

import (
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"

	time2 "github.com/henrylee2cn/ameda/test/time"
)

func TestCheckGoVersion(t *testing.T) {
	defer func() { errValueUsable = nil }()
	errValueUsable = checkGoVersion(runtime.Version())
	assert.NoError(t, errValueUsable)

	errValueUsable = checkGoVersion("go1.15")
	assert.NoError(t, errValueUsable)

	errValueUsable = checkGoVersion("go1.15rc1")
	assert.NoError(t, errValueUsable)

	errValueUsable = checkGoVersion("devel go1.19-721e999423 Mon May 30 14:13:55 2022 +0800")
	assert.NoError(t, errValueUsable)

	errValueUsable = checkGoVersion("go2.15rc1")
	assert.EqualError(t, errValueUsable, "ameda Value: required go<2.0, but current version is 'go2.15rc1'")

	errValueUsable = checkGoVersion("devel go2.19-721e999423 Mon May 30 14:13:55 2022 +0800")
	assert.EqualError(t, errValueUsable, "ameda Value: required go<2.0, but current version is 'devel go2.19-721e999423 Mon May 30 14:13:55 2022 +0800'")
}

func TestRuntimeTypeID(t *testing.T) {
	arrayEqual(t,
		RuntimeTypeIDOf(time.Time{}), RuntimeTypeID(reflect.TypeOf(time.Now())), ValueOf(time.Now()).RuntimeTypeID(), ValueFrom(reflect.ValueOf(time.Now())).RuntimeTypeID(),
	)
	arrayEqual(t,
		RuntimeTypeIDOf(&time.Time{}), RuntimeTypeID(reflect.TypeOf(&time.Time{})), ValueOf(&time.Time{}).RuntimeTypeID(), ValueFrom(reflect.ValueOf(&time.Time{})).RuntimeTypeID(),
	)
	arrayEqual(t,
		RuntimeTypeIDOf(time2.Time{}), RuntimeTypeID(reflect.TypeOf(time2.Time{S: 2})), ValueOf(time2.Time{S: 3}).RuntimeTypeID(), ValueFrom(reflect.ValueOf(time2.Time{S: 4})).RuntimeTypeID(),
	)
	arrayEqual(t,
		RuntimeTypeIDOf(&time2.Time{}), RuntimeTypeID(reflect.TypeOf(&time2.Time{S: 2})), ValueOf(&time2.Time{S: 3}).RuntimeTypeID(), ValueFrom(reflect.ValueOf(&time2.Time{S: 4})).RuntimeTypeID(),
	)
	arrayNotEqual(t, RuntimeTypeIDOf(time.Time{}), RuntimeTypeIDOf(&time.Time{}), RuntimeTypeIDOf(time2.Time{}), RuntimeTypeIDOf(&time2.Time{}))
}

func arrayEqual(t assert.TestingT, expected interface{}, actual ...interface{}) {
	if len(actual) == 0 {
		actual = append(actual, nil)
	}
	for i, a := range actual {
		assert.Equal(t, expected, a, i)
	}
}
func arrayNotEqual(t assert.TestingT, values ...interface{}) {
	if len(values) <= 1 {
		return
	}
	for i, a := range values {
		for ii, aa := range values[i+1:] {
			assert.NotEqual(t, a, aa, []int{i, ii})
		}
	}
}

func TestRuntimeTypeIDOf(t *testing.T) {
	type T1 struct {
		_ int
	}
	tid := RuntimeTypeIDOf(new(T1))
	t.Log(tid)
	assert.Equal(t, RuntimeTypeID(reflect.TypeOf(new(T1))), tid)
	tid2 := RuntimeTypeIDOf(T1{})
	assert.NotEqual(t, tid, tid2)
}

func TestKind(t *testing.T) {
	type X struct {
		A int16
		B string
	}
	var x X
	if ValueOf(&x).Kind() != reflect.Ptr {
		t.FailNow()
	}

	if ValueOf(&x).UnderlyingElem().Kind() != reflect.Struct {
		t.FailNow()
	}

	if ValueOf(x).Kind() != reflect.Struct {
		t.FailNow()
	}
	if ValueOf(x).UnderlyingElem().Kind() != reflect.Struct {
		t.FailNow()
	}

	f := func() {}
	if ValueOf(f).Kind() != reflect.Func {
		t.FailNow()
	}

	if ValueOf(t.Name).Kind() != reflect.Func {
		t.FailNow()
	}
	if ValueOf(nil).Kind() != reflect.Invalid {
		t.FailNow()
	}
	if ValueOf(nil).CanAddr() {
		t.FailNow()
	}
	if ValueOf(x).CanAddr() {
		t.FailNow()
	}
	if ValueOf(&x).CanAddr() {
		t.FailNow()
	}
	assert.True(t, ValueOf(&x).Elem().Elem().CanAddr())
	assert.True(t, ValueFrom(reflect.ValueOf(&x)).Elem().CanAddr())
	assert.True(t, ValueFrom(reflect.ValueOf(&x).Elem()).CanAddr())
}

func TestPointer(t *testing.T) {
	type X struct {
		A int16
		B string
	}
	x := X{A: 12345, B: "test"}
	if ValueOf(&x).Pointer() != reflect.ValueOf(&x).Pointer() {
		t.FailNow()
	}
	elemPtr := ValueOf(x).Pointer()
	a := *(*int16)(unsafe.Pointer(elemPtr))
	if a != x.A {
		t.FailNow()
	}
	b := *(*string)(unsafe.Pointer(elemPtr + unsafe.Offsetof(x.B)))
	if b != x.B {
		t.FailNow()
	}

	s := []string{""}
	if ValueOf(s).Pointer() != reflect.ValueOf(s).Pointer() {
		t.FailNow()
	}

	f := func() bool { return true }
	prt := ValueOf(f).Pointer()
	f = *(*func() bool)(unsafe.Pointer(&prt))
	if !f() {
		t.FailNow()
	}
	t.Log(ValueOf(f).FuncForPC().Name())
	prt = ValueOf(t.Name).Pointer()
	tName := *(*func() string)(unsafe.Pointer(&prt))
	if tName() != "TestPointer" {
		t.FailNow()
	}
	t.Log(ValueOf(t.Name).FuncForPC().Name())
	t.Log(ValueOf(s).FuncForPC() == nil)

}

func TestElem(t *testing.T) {
	type I interface{}
	var i I
	u := ValueFrom(reflect.ValueOf(i))
	type X struct {
		A int16
		B string
	}
	x := &X{A: 12345, B: "test"}
	xx := &x
	var elemPtr uintptr
	for _, v := range []interface{}{&xx, xx, x, *x} {
		val := ValueOf(v).UnderlyingElem()
		elemPtr = val.Pointer()
		a := *(*int16)(unsafe.Pointer(elemPtr))
		if !assert.Equal(t, x.A, a) {
			t.FailNow()
		}
		b := *(*string)(unsafe.Pointer(elemPtr + unsafe.Offsetof(x.B)))
		if !assert.Equal(t, x.B, b) {
			t.FailNow()
		}
	}

	var y *X
	u = ValueOf(&y)
	if !assert.False(t, u.IsNil()) {
		t.FailNow()
	}
	u = u.UnderlyingElem()
	if !assert.Equal(t, reflect.Struct, u.Kind()) {
		t.FailNow()
	}
	if !assert.True(t, u.IsNil()) {
		t.FailNow()
	}
}

func TestEmptyStruct(t *testing.T) {
	type P1 struct {
		A *int
	}
	u := ValueOf(P1{})
	if u.Pointer() != 0 {
		t.FailNow()
	}
	if !u.IsNil() {
		t.FailNow()
	}

	type P2 struct {
		A *int
		B *int
	}
	u = ValueOf(P2{})
	if u.Pointer() == 0 {
		t.FailNow()
	}
	if u.IsNil() {
		t.FailNow()
	}
}

func TestValueFrom(t *testing.T) {
	type X struct {
		A int16
		B string
	}
	x := &X{A: 12345, B: "test"}
	v := reflect.ValueOf(&x)
	u := ValueFrom2(&v).Elem()
	v = v.Elem()
	if u.RuntimeTypeID() != RuntimeTypeID(v.Type()) {
		t.FailNow()
	}
	elemPtr := u.Pointer()
	a := *(*int16)(unsafe.Pointer(elemPtr))
	if a != x.A {
		t.FailNow()
	}
	b := *(*string)(unsafe.Pointer(elemPtr + unsafe.Offsetof(x.B)))
	if b != x.B {
		t.FailNow()
	}
	if u.Pointer() != reflect.ValueOf(x).Pointer() {
		t.FailNow()
	}
}

func Benchmark_ameda(b *testing.B) {
	type T struct {
		a int
	}
	t := new(T)
	b.ReportAllocs()
	b.ResetTimer()
	u := ValueOf(t).Elem()
	for i := 0; i < b.N; i++ {
		_ = u.RuntimeTypeID()
	}
}

func Benchmark_reflect(b *testing.B) {
	type T struct {
		a int
	}
	t := new(T)
	b.ReportAllocs()
	b.ResetTimer()
	u := reflect.TypeOf(t).Elem()
	for i := 0; i < b.N; i++ {
		_ = u.String()
	}
}
