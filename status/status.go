// Package status is a handling status, similar to an error info.
package status

import (
	"encoding/json"
	"strconv"
	"unsafe"

	"github.com/henrylee2cn/goutil"
	"github.com/tidwall/gjson"
)

// Status a handling status, similar to an error info
// NOTE:
//  code: 0 means no error, -1 means unknown error
type Status struct {
	code    int32
	message string
	reason  string
}

const (
	// UnknownError unknown error status code
	UnknownError = -1
	// OK no error status code
	OK = 0
)

var (
	_ json.Marshaler   = new(Status)
	_ json.Unmarshaler = new(Status)

	reA = []byte(`{"code":`)
	reB = []byte(`,"message":`)
	reC = []byte(`,"reason":`)

	statUnknown = New(UnknownError, "Unknown Error", "")
)

// New creates a *Status.
// NOTE:
//  code=0 means no error, code=-1 means unknown error
func New(code int32, message, reason string) *Status {
	return &Status{
		code:    code,
		message: message,
		reason:  reason,
	}
}

// To converts error to *Status
// NOTE:
//  code must be -1, means unknown error status
func To(err error) *Status {
	if err == nil {
		return nil
	}
	r, ok := err.(*serror)
	if ok {
		return r.toStatus()
	}
	stat := statUnknown.Copy().SetReason(err.Error())
	return stat
}

// Copy returns the copy of Status
func (s Status) Copy() *Status {
	return &s
}

// IsOK returns true if there is no error.
func (s *Status) IsOK() bool {
	return s == nil || s.code == OK
}

// Code returns the status code.
// 0 means no error, -1 means unknown error
func (s *Status) Code() int32 {
	return s.code
}

// SetCode sets the status code.
func (s *Status) SetCode(code int32) *Status {
	s.code = code
	return s
}

// Message returns the status message displayed to the user (optional).
func (s *Status) Message() string {
	return s.message
}

// SetMessage sets the status message displayed to the user.
func (s *Status) SetMessage(message string) *Status {
	s.message = message
	return s
}

// Reason returns the cause of the status for debugging (optional).
func (s *Status) Reason() string {
	return s.reason
}

// SetReason sets the cause of the status for debugging.
func (s *Status) SetReason(reason string) *Status {
	s.reason = reason
	return s
}

// String prints status info.
func (s *Status) String() string {
	if s == nil {
		return "<nil>"
	}
	b, _ := s.MarshalJSON()
	return goutil.BytesToString(b)
}

// MarshalJSON marshals Status into JSON, implements json.Marshaler interface.
func (s *Status) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte{}, nil
	}
	var b = append(reA, strconv.FormatInt(int64(s.code), 10)...)
	if len(s.message) > 0 {
		b = append(b, reB...)
		b = append(b, goutil.StringMarshalJSON(s.message, false)...)
	}
	if len(s.reason) > 0 {
		b = append(b, reC...)
		b = append(b, goutil.StringMarshalJSON(s.reason, false)...)
	}
	b = append(b, '}')
	return b, nil
}

// UnmarshalJSON unmarshals a JSON description of self.
func (s *Status) UnmarshalJSON(b []byte) error {
	if s == nil {
		return nil
	}
	str := goutil.BytesToString(b)
	s.code = int32(gjson.Get(str, "code").Int())
	s.message = gjson.Get(str, "message").String()
	s.reason = gjson.Get(str, "reason").String()
	return nil
}

// ToError converts to error interface.
func (s *Status) ToError() error {
	if s == nil || s.code == OK {
		return nil
	}
	return (*serror)(unsafe.Pointer(s))
}

type serror Status

func (r *serror) Error() string {
	b, _ := r.toStatus().MarshalJSON()
	return goutil.BytesToString(b)
}

func (r *serror) toStatus() *Status {
	return (*Status)(unsafe.Pointer(r))
}
