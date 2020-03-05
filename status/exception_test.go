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

func TestCatchNil20(t *testing.T) {
	var stat *Status
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	Panic(nil, true)
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

func TestCatchNil22(t *testing.T) {
	var stat = New(1, "text")
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	Panic(stat, true)
}

func TestCatchNil23(t *testing.T) {
	var stat *Status
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	stat.Panic(false)
}

func TestCatchNil24(t *testing.T) {
	var stat *Status
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	stat.Panic(true)
}

func TestCatchNil25(t *testing.T) {
	var stat = NewWithStack(1, "text")
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	stat.Panic(false)
}

func TestCatchNil26(t *testing.T) {
	var stat = NewWithStack(1, "text")
	var realStat bool
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.True(t, realStat)
	}()
	defer Catch(&stat, &realStat)
	stat.Panic(true)
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
	stat.CopyCheck(err)
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
	stat.CopyThrow("a test error")
}

func TestThrow3(t *testing.T) {
	var stat = New(400, "bad", "raw cause")
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
		assert.Equal(t, "raw cause", stat.cause.Error())
	}()
	defer Catch(&stat)
	stat.CopyThrow()
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
