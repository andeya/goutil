// Package status is a handling status with code, msg, cause and stack.
package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/henrylee2cn/goutil"
)

// Status a handling status with code, msg, cause and stack.
type Status struct {
	code  int32
	msg   string
	cause error
	*stack
}

// OK status
const OK int32 = 0

// New creates a handling status with code, msg and cause.
// NOTE:
//  code=0 means no error
func New(code int32, msg string, cause error) *Status {
	return &Status{
		code:  code,
		msg:   msg,
		cause: cause,
	}
}

// WithStack creates a handling status with code, msg, cause and stack.
// NOTE:
//  code=0 means no error
func WithStack(code int32, msg string, cause error) *Status {
	s := New(code, msg, cause)
	s.stack = callers()
	return s
}

// Copy returns the copy of Status.
func (s *Status) Copy(withStack bool, newCause ...error) *Status {
	if s == nil {
		return nil
	}
	copy := *s
	if withStack {
		copy.stack = callers()
	} else {
		copy.stack = nil
	}
	if len(newCause) > 0 {
		copy.cause = newCause[0]
	}
	return &copy
}

// Code returns the status code.
func (s *Status) Code() int32 {
	if s == nil {
		return OK
	}
	return s.code
}

// Msg returns the status msg displayed to the user (optional).
func (s *Status) Msg() string {
	if s == nil {
		return ""
	}
	return s.msg
}

// Cause returns the cause of the status for debugging (optional).
func (s *Status) Cause() error {
	if s == nil {
		return nil
	}
	return s.cause
}

// StackTrace returns stack trace.
func (s *Status) StackTrace() StackTrace {
	if s.stack == nil {
		return nil
	}
	return s.stack.StackTrace()
}

// String prints status info.
func (s *Status) String() string {
	if s == nil {
		return "<nil>"
	}
	b, _ := s.MarshalJSON()
	return goutil.BytesToString(b)
}

// Format implementes fmt.Formatter.
func (s *Status) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			fmt.Fprintf(state, "%+v", s.String())
			if s.stack != nil {
				s.stack.Format(state, verb)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(state, s.String())
	case 'q':
		fmt.Fprintf(state, "%q", s.String())
	}
}

var (
	_ json.Marshaler   = new(Status)
	_ json.Unmarshaler = new(Status)

	reA = []byte(`{"code":`)
	reB = []byte(`,"msg":`)
	reC = []byte(`,"cause":`)
)

// MarshalJSON marshals Status into JSON, implements json.Marshaler interface.
func (s *Status) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte{}, nil
	}
	var b = append(reA, strconv.FormatInt(int64(s.code), 10)...)
	if len(s.msg) > 0 {
		b = append(b, reB...)
		b = append(b, goutil.StringMarshalJSON(s.msg, false)...)
	}
	if s.cause != nil {
		if cause := s.cause.Error(); cause != "" {
			b = append(b, reC...)
			b = append(b, goutil.StringMarshalJSON(cause, false)...)
		}
	}
	b = append(b, '}')
	return b, nil
}

// UnmarshalJSON unmarshals a JSON description of self.
func (s *Status) UnmarshalJSON(b []byte) error {
	if s == nil {
		return nil
	}
	var v = &struct {
		Code  int32  `json:"code"`
		Msg   string `json:"msg"`
		Cause string `json:"cause"`
	}{}
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	s.code = v.Code
	s.msg = v.Msg
	if v.Cause != "" {
		s.cause = errors.New(v.Cause)
	}
	return nil
}
