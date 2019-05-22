package status

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatchNil(t *testing.T) {
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.Empty(t, stat)
	}()
	defer Catch(&stat)
	panic(nil)
}

func TestCatchNil2(t *testing.T) {
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
	}()
	defer Catch(&stat)
	Panic(nil)
}

func TestCatchNotNil(t *testing.T) {
	var stat *Status
	defer func() {
		t.Logf("%+v", stat)
		assert.True(t, stat != nil)
	}()

	defer Catch(&stat)
	Check(errors.New("a test error"), 400, "check")
}

func TestThrow(t *testing.T) {
	defer func() {
		r := recover()
		t.Logf("%+v", r)
		assert.True(t, r != nil)
	}()

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
