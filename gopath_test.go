package goutil

import (
	"os"
	"testing"
)

func TestGetFirstGopath(t *testing.T) {
	os.Setenv("GOPATH", "")
	targetGoPath := os.Args[0] + string(os.PathSeparator)
	os.Args[0] = os.Args[0] + "/src/"
	goPath, err := GetFirstGopath(false)
	if err == nil {
		t.FailNow()
	}
	goPath, err = GetFirstGopath(true)
	if err != nil {
		t.FailNow()
	}
	if goPath != targetGoPath {
		t.Fatalf("expect %s, but get %s", targetGoPath, goPath)
	}
	os.Setenv("GOPATH", targetGoPath)
	goPath, err = GetFirstGopath(false)
	if err != nil {
		t.FailNow()
	}
	if goPath != targetGoPath {
		t.Fatalf("expect %s, but get %s", targetGoPath, goPath)
	}

}
