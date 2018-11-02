package goutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetFirstGopath gets the first $GOPATH value.
func GetFirstGopath(allowAutomaticGuessing bool) (goPath string, err error) {
	goPath = os.Getenv("GOPATH")
	defer func() {
		goPath = strings.Replace(goPath, "/", string(os.PathSeparator), -1)
	}()
	if len(goPath) == 0 {
		if !allowAutomaticGuessing {
			err = errors.New("not found GOPATH")
			return
		}
		p, _ := filepath.Abs(os.Args[0])
		p = strings.Replace(p, "\\", "/", -1) + "/"
		i := strings.Index(p, "/src/")
		if i == -1 {
			fmt.Println(p)
			err = errors.New("not found GOPATH")
			return
		}
		goPath = p[:i+1]
		return
	}
	var sep string
	if runtime.GOOS == "windows" {
		sep = ";"
	} else {
		sep = ":"
	}
	if goPaths := strings.Split(goPath, sep); len(goPaths) > 1 {
		goPath = goPaths[0]
	}
	goPath, _ = filepath.Abs(goPath)
	goPath = strings.Replace(goPath, "\\", "/", -1)
	if goPath[len(goPath)-1] != '/' {
		goPath += "/"
	}
	return
}
