// errors is improved errors package.
package errors

import (
	"fmt"

	"github.com/henrylee2cn/goutil"
)

// New returns an error that formats as the given text.
func New(text string) error {
	return &myerror{text}
}

// myerror is a trivial implementation of error.
type myerror struct {
	s string
}

func (e *myerror) Error() string {
	return e.s
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(format string, a ...interface{}) error {
	return New(fmt.Sprintf(format, a...))
}

// Separator multi errors separator
const Separator byte = 0xA

// Merge merges multi errors.
func Merge(errs ...error) error {
	if len(errs) == 1 {
		return errs[0]
	}
	var b = make([]byte, 0, 32)
	var e string
	for _, err := range errs {
		if err != nil {
			e = err.Error()
			if e[0] != Separator {
				b = append(b, Separator)
			}
			b = append(b, e...)
		}
	}
	if len(b) == 0 {
		return nil
	}
	return New(goutil.BytesToString(b))
}
