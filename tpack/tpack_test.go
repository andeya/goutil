package tpack

import (
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestRuntimeTypeID(t *testing.T) {
	type (
		GoTime = time.Time
		Time   time.Time
		I2     interface {
			String() string
		}
		I1 interface {
			UnixNano() int64
			I2
		}
	)
	t0 := new(time.Time)
	t1 := Unpack(t0).RuntimeTypeID()
	t2 := Unpack(new(GoTime)).RuntimeTypeID()
	t3 := Unpack(new(Time)).RuntimeTypeID()
	t.Log(t1, t2, t3)
	e0 := time.Time{}
	e1 := Unpack(e0).RuntimeTypeID()
	e2 := Unpack(GoTime{}).RuntimeTypeID()
	e3 := Unpack(Time{}).RuntimeTypeID()
	i := Unpack(I2(I1(&GoTime{}))).RuntimeTypeID()
	if t1 != t2 || t1 != e1 || t1 != e2 || t1 != i || t3 != e3 {
		t.FailNow()
	}
	t.Log(e1, e2, e3, i, RuntimeTypeID(reflect.TypeOf(t0)), Unpack(t0.String).RuntimeTypeID())
}

func TestKind(t *testing.T) {
	type X struct {
		A int16
		B string
	}
	var x X
	if Unpack(&x).Kind() != reflect.Ptr {
		t.FailNow()
	}

	if Unpack(&x).UnderlyingElem().Kind() != reflect.Struct {
		t.FailNow()
	}

	if Unpack(x).Kind() != reflect.Struct {
		t.FailNow()
	}
	if Unpack(x).UnderlyingElem().Kind() != reflect.Struct {
		t.FailNow()
	}

	f := func() {}
	if Unpack(f).Kind() != reflect.Func {
		t.FailNow()
	}

	if Unpack(t.Name).Kind() != reflect.Func {
		t.FailNow()
	}
}

func TestPointer(t *testing.T) {
	type X struct {
		A int16
		B string
	}
	x := X{A: 12345, B: "test"}
	if Unpack(&x).Pointer() != reflect.ValueOf(&x).Pointer() {
		t.FailNow()
	}
	elemPtr := Unpack(x).Pointer()
	a := *(*int16)(unsafe.Pointer(elemPtr))
	if a != x.A {
		t.FailNow()
	}
	b := *(*string)(unsafe.Pointer(elemPtr + unsafe.Offsetof(x.B)))
	if b != x.B {
		t.FailNow()
	}

	s := []string{""}
	if Unpack(s).Pointer() != reflect.ValueOf(s).Pointer() {
		t.FailNow()
	}

	f := func() bool { return true }
	prt := Unpack(f).Pointer()
	f = *(*func() bool)(unsafe.Pointer(&prt))
	if !f() {
		t.FailNow()
	}
	t.Log(Unpack(f).FuncForPC().Name())
	prt = Unpack(t.Name).Pointer()
	tName := *(*func() string)(unsafe.Pointer(&prt))
	if tName() != "TestPointer" {
		t.FailNow()
	}
	t.Log(Unpack(t.Name).FuncForPC().Name())
	t.Log(Unpack(s).FuncForPC() == nil)

}

func TestElem(t *testing.T) {
	type I interface{}
	var i I
	u := From(reflect.ValueOf(i))
	u.Elem()

	type X struct {
		A int16
		B string
	}
	x := &X{A: 12345, B: "test"}
	xx := &x
	var elemPtr uintptr
	for i, v := range []interface{}{&xx, xx, x, *x} {
		if i == 0 {
			elemPtr = Unpack(v).UnderlyingElem().Pointer()
		} else {
			elemPtr = Unpack(v).Elem().Pointer()
		}
		a := *(*int16)(unsafe.Pointer(elemPtr))
		if a != x.A {
			t.FailNow()
		}
		b := *(*string)(unsafe.Pointer(elemPtr + unsafe.Offsetof(x.B)))
		if b != x.B {
			t.FailNow()
		}
	}

	var y *X
	u = Unpack(&y)
	if u.IsNil() {
		t.FailNow()
	}
	u = u.UnderlyingElem()
	if u.Kind() != reflect.Struct {
		t.FailNow()
	}
	if !u.IsNil() {
		t.FailNow()
	}
}

func TestEmptyStruct(t *testing.T) {
	type P1 struct {
		A *int
	}
	u := Unpack(P1{})
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
	u = Unpack(P2{})
	if u.Pointer() == 0 {
		t.FailNow()
	}
	if u.IsNil() {
		t.FailNow()
	}
}

func TestFrom(t *testing.T) {
	type X struct {
		A int16
		B string
	}
	x := &X{A: 12345, B: "test"}
	v := reflect.ValueOf(&x)
	u := From(v).Elem()
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

func Benchmark_tpack(b *testing.B) {
	type T struct {
		a int
	}
	t := new(T)
	b.ReportAllocs()
	b.ResetTimer()
	u := Unpack(t).Elem()
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
