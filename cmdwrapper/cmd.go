// Package cmdwrapper exec cmd and catch the log the result.
package cmdwrapper

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"unsafe"
)

var cmdArg [2]string

func init() {
	if runtime.GOOS == "windows" {
		cmdArg[0] = "cmd"
		cmdArg[1] = "/c"
	} else {
		cmdArg[0] = "/bin/sh"
		cmdArg[1] = "-c"
	}
}

// Run exec cmd and catch the log the result.
func Run(cmdLine string) *Result {
	cmd := exec.Command(cmdArg[0], cmdArg[1], cmdLine)
	var ret = new(Result)
	cmd.Stdout = &ret.buf
	cmd.Stderr = &ret.buf
	cmd.Env = os.Environ()
	ret.err = cmd.Run()
	return ret
}

// Result cmd exec result
type Result struct {
	buf bytes.Buffer
	err error
	str *string
}

// Err returns the error log.
func (r *Result) Err() error {
	if r.err == nil {
		return nil
	}
	r.err = errors.New(r.String())
	return r.err
}

// String returns the exec log.
func (r *Result) String() string {
	if r.str == nil {
		b := append(bytes.TrimSpace(r.buf.Bytes()), ' ', '(')
		b = append(b, r.err.Error()...)
		b = append(b, ')')
		r.str = (*string)(unsafe.Pointer(&b))
	}
	return *r.str
}
