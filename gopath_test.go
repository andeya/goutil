package goutil

import (
	"os"
	"strings"
	"testing"
)

func TestGetFirstGopath(t *testing.T) {
	targetGoPath := os.Getenv("GOPATH")
	targetGoPath = strings.SplitN(strings.SplitN(targetGoPath, ":", 1)[0], ";", 1)[0]
	targetGoPath = strings.TrimRight(strings.TrimRight(targetGoPath, "/"), "\\") + string(os.PathSeparator)
	os.Setenv("GOPATH", "")
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
