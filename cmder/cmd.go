// Package cmder exec cmd and catch the result.
package cmder

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
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

// Run exec cmd and catch the result.
// Waits for the given command to finish with a timeout.
// If the command times out, it attempts to kill the process.
func Run(cmdLine string, timeout ...time.Duration) *Result {
	cmd := exec.Command(cmdArg[0], cmdArg[1], cmdLine)
	var ret = new(Result)
	cmd.Stdout = &ret.buf
	cmd.Stderr = &ret.buf
	cmd.Env = os.Environ()
	ret.err = cmd.Start()
	if ret.err != nil {
		return ret
	}
	if len(timeout) == 0 || timeout[0] <= 0 {
		ret.err = cmd.Wait()
		return ret
	}
	timer := time.NewTimer(timeout[0])
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case ret.err = <-done:
		timer.Stop()
	case <-timer.C:
		if err := cmd.Process.Kill(); err != nil {
			ret.err = fmt.Errorf("command timed out and killing process fail: %s", err.Error())
		} else {
			// wait for the command to return after killing it
			<-done
			ret.err = errors.New("command timed out")
		}
	}
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
		b := bytes.TrimSpace(r.buf.Bytes())
		if r.err != nil {
			b = append(b, ' ', '(')
			b = append(b, r.err.Error()...)
			b = append(b, ')')
		}
		r.str = (*string)(unsafe.Pointer(&b))
	}
	return *r.str
}
