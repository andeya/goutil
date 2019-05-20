package status

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestStatus(t *testing.T) {
	stat := New(
		400,
		"msg...",
		"bala...bala...",
	)
	t.Logf("%v", stat)

	err := errors.New("xxxxxxxxxx")
	stat = testWithStack(err)
	t.Logf("%+v", stat)

	stat = new(Status)
	err = json.Unmarshal([]byte(`{"code":404,"msg":"Not Found","cause":"xxxxxxxxxx"}`), stat)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", stat)
	t.Logf("%+v", stat.Copy(true, "zzzzzzzzz"))
}

func testWithStack(err error) *Status {
	return NewWithStack(404, "Not Found", err)
}
