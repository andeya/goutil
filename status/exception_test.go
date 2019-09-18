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

func TestCatchNil3(t *testing.T) {
	var realStat bool
	defer func() {
		assert.False(t, realStat)
	}()
	defer Catch(nil, &realStat)
}

func TestCatchStatus(t *testing.T) {
	var stat = NewWithStack(2, "TestCatchStatus", nil)
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

func TestThrow(t *testing.T) {
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
	}()

	defer Catch(&stat)
	Throw(400, "", "a test error")
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
