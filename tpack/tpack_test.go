package tpack

import (
	"reflect"
	"testing"
	"time"
)

func TestGetTypeID(t *testing.T) {
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
	t1 := Unpack(t0).TypeID()
	t2 := Unpack(new(GoTime)).TypeID()
	t3 := Unpack(new(Time)).TypeID()
	t.Log(t1, t2, t3)
	e0 := time.Time{}
	e1 := Unpack(e0).TypeID()
	e2 := Unpack(GoTime{}).TypeID()
	e3 := Unpack(Time{}).TypeID()
	i := Unpack(I2(I1(&GoTime{}))).TypeID()
	if t1 != t2 || t1 != e1 || t1 != e2 || t1 != i || t3 != e3 {
		t.FailNow()
	}
	t.Log(e1, e2, e3, i, TypeID(reflect.TypeOf(t0)))
}

func BenchmarkUnpack_pointer(b *testing.B) {
	b.StopTimer()
	type T struct {
		a int
	}
	var t T
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = Unpack(t).TypeID()
	}
}

func BenchmarkValueOf_go(b *testing.B) {
	b.StopTimer()
	type T struct {
		a int
	}
	var t T
	v := reflect.ValueOf(t)
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = v.String()
	}
}
