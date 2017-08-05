package goutil

import (
	"path/filepath"
	"reflect"
	"testing"
)

var noExistedFile = "/tmp/not_existed_file"

func TestSelfPath(t *testing.T) {
	path := SelfPath()
	if path == "" {
		t.Error("path cannot be empty")
	}
	t.Logf("SelfPath: %s", path)
}

func TestSelfDir(t *testing.T) {
	dir := SelfDir()
	t.Logf("SelfDir: %s", dir)
}

func TestSelfChdir(t *testing.T) {
	SelfChdir()
	path, err := filepath.Abs("a")
	t.Logf("SelfChdir: %s %v", path, err)
}

func TestFileExists(t *testing.T) {
	if !FileExists("./file.go") {
		t.Errorf("./file.go should exists, but it didn't")
	}

	if FileExists(noExistedFile) {
		t.Errorf("Weird, how could this file exists: %s", noExistedFile)
	}
}

func TestSearchFile(t *testing.T) {
	path, err := SearchFile(filepath.Base(SelfPath()), SelfDir())
	if err != nil {
		t.Error(err)
	}
	t.Log(path)

	_, err = SearchFile(noExistedFile, ".")
	if err == nil {
		t.Errorf("err shouldnot be nil, got path: %s", SelfDir())
	}
}

func TestGrepFile(t *testing.T) {
	_, err := GrepFile("", noExistedFile)
	if err == nil {
		t.Error("expect file-not-existed error, but got nothing")
	}

	lines, err := GrepFile(`^func GrepFile.*{$`, "file.go")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(lines, []string{"func GrepFile(patten string, filename string) (lines []string, err error) {"}) {
		t.Errorf("expect [\"func GrepFile(patten string, filename string) (lines []string, err error) {\"], but receive %v", lines)
	}
}
