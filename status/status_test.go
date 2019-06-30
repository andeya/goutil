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
	t.Logf("%+v", stat.Copy("zzzzzzzzz", 0))
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

func TestStatusStack(t *testing.T) {
	jsonBytes := []byte(`{"code":404,"msg":"Not Found","cause":"xxxxxxxxxx"}`)
	s, err := FromJSON(jsonBytes, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", s)
	queryBytes := []byte(`code=400&msg=msg...&cause=bala...bala...`)
	s = FromQuery(queryBytes, true)
	t.Logf("%+v", s)
	s = New(400, "tag", "test").TagStack()
	t.Logf("%+v", s)
}

func TestStatusPanic(t *testing.T) {
	defer func() {
		stat := New(1, "panic", recover()).TagStack(3)
		t.Logf("%+v", stat)
	}()
	var a []byte
	_ = a[1]
}

func TestStatusPanic2(t *testing.T) {
	defer func() {
		stat := New(1, "panic", recover()).TagStack(2)
		t.Logf("%+v", stat)
	}()
	panic("this is panic text")
}

func TestStatusPanic3(t *testing.T) {
	stat := New(0, "", nil)
	defer func() {
		stat = stat.Copy(recover(), 3)
		t.Logf("%+v", stat)
	}()
	var a []byte
	_ = a[1]
}

func TestStatusPanic4(t *testing.T) {
	stat := New(0, "", nil)
	defer func() {
		stat = stat.Copy(recover(), 2)
		t.Logf("%+v", stat)
	}()
	panic("this is panic text")
}

func testWithStack(err error) *Status {
	return NewWithStack(404, "Not Found", err)
}
