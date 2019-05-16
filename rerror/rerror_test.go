package rerror

import (
	"errors"
	"testing"
)

func Test(t *testing.T) {
	var rerr *Rerror
	t.Logf("%v", rerr.HasError())
	rerr = new(Rerror)
	t.Logf("%v", rerr.HasError())
	t.Logf("%v", rerr)
	rerr.Code = 400
	rerr.Message = "msg"
	t.Logf("%v", rerr)
	rerr.Reason = `"bala...bala..."`
	t.Logf("%v", rerr)
	err := rerr.ToError()
	t.Logf("test ToError 1: %v", err)
	rerr = To(err)
	t.Logf("test To 1: %s", rerr)
	rerr = nil
	err = rerr.ToError()
	t.Logf("test ToError 2: %v", err)
	rerr = To(nil)
	t.Logf("test To 2: %s", rerr)
	rerr = To(errors.New("text error"))
	t.Logf("test To 3: %s", rerr)
}
