package status

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestStatusJSON(t *testing.T) {
	if (*Status)(nil).String() != "<nil>" {
		t.FailNow()
	}
	b, _ := json.Marshal((*Status)(nil))
	if !bytes.Equal(b, null) {
		t.FailNow()
	}
	b, _ = (*Status)(nil).MarshalJSON()
	if !bytes.Equal(b, null) {
		t.FailNow()
	}

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

func TestStatusQuery(t *testing.T) {
	if (*Status)(nil).QueryString() != "" {
		t.FailNow()
	}
	b := (*Status)(nil).EncodeQuery()
	if b != nil {
		t.FailNow()
	}

	expect := New(
		400,
		"msg...",
		"bala...bala...",
	)
	t.Logf("%v", expect.QueryString())

	stat := new(Status)
	stat.DecodeQuery([]byte(expect.QueryString()))
	if !reflect.DeepEqual(*expect, *stat) {
		t.Fatalf("got:%s, want:%s", stat, expect)
	}
}

func testWithStack(err error) *Status {
	return NewWithStack(404, "Not Found", err)
}
