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

const (
	// OK status
	OK int32 = 0

	// UnknownError status
	UnknownError int32 = -1
)

// Status a handling status with code, msg, cause and stack.
type Status struct {
	code  int32
	msg   string
	cause error
	*stack
}

// New creates a handling status with code, msg and cause.
// NOTE:
//  code=0 means no error
func New(code int32, msg string, cause interface{}) *Status {
	s := &Status{
		code: code,
		msg:  msg,
	}
	switch v := cause.(type) {
	case nil:
	case error:
		s.cause = v
	case string:
		s.cause = errors.New(v)
	case *Status:
		s.cause = v.cause
	case Status:
		s.cause = v.cause
	default:
		s.cause = fmt.Errorf("%v", v)
	}
	return s
}

// NewWithStack creates a handling status with code, msg and cause and stack.
// NOTE:
//  code=0 means no error
func NewWithStack(code int32, msg string, cause interface{}) *Status {
	s := New(code, msg, cause)
	s.stack = callers()
	return s
}

// Copy returns the copy of Status.
func (s *Status) Copy(withStack bool, newCause ...interface{}) *Status {
	if s == nil {
		return nil
	}
	var cause interface{} = s.cause
	if len(newCause) > 0 {
		cause = newCause[0]
	}
	copy := New(s.code, s.msg, cause)
	if withStack {
		copy.stack = callers()
	}
	return copy
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

// OK returns whether is OK status (code=0).
func (s *Status) OK() bool {
	return s.Code() == OK
}

// UnknownError returns whether is UnknownError status (code=-1).
func (s *Status) UnknownError() bool {
	return s.Code() == UnknownError
}

// StackTrace returns stack trace.
func (s *Status) StackTrace() StackTrace {
	if s == nil || s.stack == nil {
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

type exportStatus struct {
	Code  int32  `json:"code"`
	Msg   string `json:"msg"`
	Cause string `json:"cause"`
}

var (
	reA  = []byte(`{"code":`)
	reB  = []byte(`,"msg":`)
	reC  = []byte(`,"cause":`)
	null = []byte("null")
)
var (
	_ json.Marshaler   = new(Status)
	_ json.Unmarshaler = new(Status)
)

// MarshalJSON marshals Status into JSON, implements json.Marshaler interface.
func (s *Status) MarshalJSON() ([]byte, error) {
	if s == nil {
		return null, nil
	}
	b := append(reA, strconv.FormatInt(int64(s.code), 10)...)

	b = append(b, reB...)
	b = append(b, goutil.StringMarshalJSON(s.msg, false)...)

	var cause string
	if s.cause != nil {
		cause = s.cause.Error()
	}
	b = append(b, reC...)
	b = append(b, goutil.StringMarshalJSON(cause, false)...)

	b = append(b, '}')
	return b, nil
}

// UnmarshalJSON unmarshals a JSON description of self.
func (s *Status) UnmarshalJSON(b []byte) error {
	if s == nil {
		return nil
	}
	var v exportStatus
	err := json.Unmarshal(b, &v)
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
