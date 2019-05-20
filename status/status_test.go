package status

import (
	"encoding/json"
	"errors"
	"testing"
)

func Test(t *testing.T) {
	stat := New(
		400,
		"msg",
		errors.New("bala...bala..."),
	)
	t.Logf("%v", stat)

	stat = testWithStack()
	t.Logf("%+v", stat)

	stat = new(Status)
	err := json.Unmarshal([]byte(`{"code":404,"msg":"Not Found","cause":"xxxxxxxxxx"}`), stat)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", stat)
	t.Logf("%+v", stat.Copy(true, errors.New("zzzzzzzzz")))
}

func testWithStack() *Status {
	return WithStack(404, "Not Found", errors.New("xxxxxxxxxx"))
}
