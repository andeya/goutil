package status

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatchNil(t *testing.T) {
	var stat *Status
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.Empty(t, stat)
		assert.True(t, stat != nil)
		assert.False(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	panic(nil)
}

func TestCatchNil10(t *testing.T) {
	var stat *Status
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.Empty(t, stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	panic(stat)
}

func TestCatchNil11(t *testing.T) {
	var stat *Status
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer CatchWithStack(&stat, &realStat)
	panic(stat)
}

func TestCatchNil2(t *testing.T) {
	var stat *Status
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	Panic(nil)
}

func TestCatchNil21(t *testing.T) {
	var stat = New(1, "text")
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	Panic(stat)
}

func TestCatchNil3(t *testing.T) {
	var realStat bool
	defer func() {
		assert.False(t, realStat)
	}()
	defer Catch(nil, &realStat)
}

func TestCatchStatus(t *testing.T) {
	var stat = NewWithStack(2, "TestCatchStatus")
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.Equal(t, int32(2), stat.Code())
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
}

func TestCatchNotNil(t *testing.T) {
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
	}()

	defer Catch(&stat)
	Check(errors.New("a test error"), 400, "check 1")
}

func TestCheckWhenError(t *testing.T) {
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
	}()

	defer Catch(&stat)
	Check(errors.New("a test error"), 400, "check 2", func() { t.Log("whenerror") })
}

func TestCheckError(t *testing.T) {
	err := errors.New("a test error 3")
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.Equal(t, err.Error(), stat.Msg())
	}()
	defer Catch(&stat)
	Check(err, 400, "")
}

func TestCheckError2(t *testing.T) {
	err := errors.New("a test error 3")
	var stat = New(400, "bad", "raw cause")
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.Equal(t, err, stat.Cause())
	}()
	defer Catch(&stat)
	stat.NewCheck(err)
}

func TestThrow(t *testing.T) {
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
	}()
	defer Catch(&stat)
	Throw(400, "", "a test error")
}

func TestThrow2(t *testing.T) {
	var stat = New(400, "bad", "raw cause")
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.Equal(t, "a test error", stat.Cause().Error())
	}()
	defer Catch(&stat)
	stat.NewThrow("a test error")
}

func TestThrow3(t *testing.T) {
	var stat = New(400, "bad", "raw cause")
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.Equal(t, "raw cause", stat.cause.Error())
	}()
	defer Catch(&stat)
	stat.NewThrow()
}

func TestCatchNotError(t *testing.T) {
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
	}()

	defer Catch(&stat)
	panic("this is not a error")
}

func TestCatchNotError2(t *testing.T) {
	var stat = New(400, "", "old error")
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
	}()

	defer Catch(&stat)
	panic("this is not a error")
}

func TestFindPanicStack(t *testing.T) {
	defer func() {
		recover()
		stack := findPanicStack()
		t.Logf("%+v", stack)
	}()
	panic("this is not a error")
}

func TestPanicStackTrace(t *testing.T) {
	defer func() {
		stack := PanicStackTrace()
		t.Logf("%+v", stack)
	}()
}

func TestGetStackTrace(t *testing.T) {
	defer func() {
		stack := GetStackTrace(0)
		t.Logf("%+v", stack)
	}()
}
