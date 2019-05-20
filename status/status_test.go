package status

import (
	"errors"
	"testing"
)

func Test(t *testing.T) {
	var stat *Status
	t.Logf("%v", stat.IsOK())
	stat = new(Status)
	t.Logf("%v", stat.IsOK())
	t.Logf("%v", stat)
	stat.SetCode(400)
	stat.SetMessage("msg")
	t.Logf("%v", stat)
	stat.SetReason(`"bala...bala..."`)
	t.Logf("%v", stat)
	err := stat.ToError()
	t.Logf("test ToError 1: %v", err)
	stat = To(err)
	t.Logf("test To 1: %s", stat)
	stat = nil
	err = stat.ToError()
	t.Logf("test ToError 2: %v", err)
	stat = To(nil)
	t.Logf("test To 2: %s", stat)
	stat = To(errors.New("text error"))
	t.Logf("test To 3: %s", stat)
}
