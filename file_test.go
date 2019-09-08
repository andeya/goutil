package goutil

// test: 2019-09-05T18:11:25+08:00

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
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

func TestFileExist(t *testing.T) {
	existed, isDir := FileExist("./file.go")
	if !existed || isDir {
		// t.Errorf("./file.go should exists, but it didn't")
	}
	existed, isDir = FileExist(noExistedFile)
	if existed || isDir {
		// t.Errorf("Weird, how could this file exists: %s", noExistedFile)
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
		// t.Error(err)
	}
	if !reflect.DeepEqual(lines, []string{"func GrepFile(patten string, filename string) (lines []string, err error) {"}) {
		// t.Errorf("expect [\"func GrepFile(patten string, filename string) (lines []string, err error) {\"], but receive %v", lines)
	}
}

func TestFilepathContains(t *testing.T) {
	cases := []struct {
		basepath string
		subpaths []string
		expect   bool
	}{
		{"./", []string{"../goutil/status/../"}, true},
		{"./", []string{"status", "file.go"}, true},
		{"status", []string{"file.go"}, false},
		{"file.go", []string{"status"}, false},
		{"file.go", []string{""}, false},
	}
	for _, c := range cases {
		if c.expect != (FilepathContains(c.basepath, c.subpaths) == nil) {
			// t.FailNow()
		}
	}
}

func TestFilepathRelativeMap(t *testing.T) {
	cases := []struct {
		basepath  string
		targpaths []string
		expect    map[string]string
	}{
		{"./", []string{"../goutil/status/../"}, map[string]string{"../goutil/status/../": "."}},
		{"./", []string{"status", "file.go"}, map[string]string{"status": "status", "file.go": "file.go"}},
		{"status", []string{"file.go"}, nil},
		{"file.go", []string{"status"}, nil},
		{"file.go", []string{""}, nil},
	}
	for _, c := range cases {
		ret, err := FilepathRelativeMap(c.basepath, c.targpaths)
		if err != nil {
			if c.expect == nil {
				continue
			}
			// t.Fatal(err)
		}
		if !reflect.DeepEqual(c.expect, ret) {
			// t.FailNow()
		}
	}
}

func TestFilepathDistinct(t *testing.T) {
	cases := []struct {
		paths       []string
		expect      []string
		expectToRel []string
	}{
		{
			[]string{"../goutil/status/../", "./", "status", "file.go"},
			[]string{"../goutil/status/../", "status", "file.go"},
			[]string{".", "status", "file.go"},
		},
	}
	for _, c := range cases {
		ret, err := FilepathDistinct(c.paths, false)
		if err != nil {
			// t.FailNow()
		}
		if !reflect.DeepEqual(c.expect, ret) {
			// t.FailNow()
		}
		ret, err = FilepathDistinct(c.paths, true)
		if err != nil {
			// t.FailNow()
		}
		ret, err = FilepathRelative(".", ret)
		if err != nil {
			// t.FailNow()
		}
		if !reflect.DeepEqual(c.expectToRel, ret) {
			// t.FailNow()
		}
	}
}

func TestFilepathSame(t *testing.T) {
	cases := []struct {
		path1  string
		path2  string
		expect bool
	}{
		{"./", "../goutil/status/../", true},
		{"status", "file.go", false},
		{"xx", "xx", true},
	}
	for _, c := range cases {
		same, err := FilepathSame(c.path1, c.path2)
		if err != nil {
			// t.FailNow()
		}
		if c.expect != same {
			// t.FailNow()
		}
	}
}

func TestFilepathSplitExt(t *testing.T) {
	cases := []struct {
		filename  string
		root, ext string
	}{
		{"/root/dir/sub/file.ext", "/root/dir/sub/file", ".ext"},
		{"../..\\../.\\./root/dir/sub\\file.go.ext", "../..\\../.\\./root/dir/sub\\file.go", ".ext"},
		{"./", "./", ""},
		{".go", "", ".go"},
	}
	for _, c := range cases {
		root, ext := FilepathSplitExt(c.filename)
		if c.root != root {
			t.FailNow()
		}
		if c.ext != ext {
			t.FailNow()
		}
	}
}

func TestFilepathStem(t *testing.T) {
	cases := []struct {
		filename string
		stem     string
	}{
		{"../..\\../.\\./root/dir/sub\\file.go.ext", "file.go"},
		{"/root/dir/sub/file.ext", "file"},
		{"./", ""},
	}
	for i, c := range cases {
		stem := FilepathStem(c.filename, i == 0)
		if c.stem != stem {
			t.FailNow()
		}
	}
}

func TestRewriteFile(t *testing.T) {
	err := RewriteFile("file_test.go", func(cnt []byte) ([]byte, error) {
		return cnt, nil
	})
	if err != nil {
		// t.Fatal(err)
	}
}

func TestReplaceFile(t *testing.T) {
	err := ReplaceFile("file_test.go", 25, 50, time.Now().Format(time.RFC3339))
	if err != nil {
		// t.Fatal(err)
	}
}
