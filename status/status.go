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
type Status struct {
	// Code status code
	// 0 means no error, -1 means unknown error
	Code int32 `json:"code"`
	// Message the status message displayed to the user (optional)
	Message string `json:"message"`
	// Reason the cause of the status for debugging (optional)
	Reason string `json:"reason"`
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
		Code:    code,
		Message: message,
		Reason:  reason,
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
func (r Status) Copy() *Status {
	return &r
}

// IsOK returns true if there is no error.
func (r *Status) IsOK() bool {
	return r == nil || r.Code == OK
}

// SetMessage sets the status message displayed to the user.
func (r *Status) SetMessage(message string) *Status {
	r.Message = message
	return r
}

// SetReason sets the cause of the status for debugging.
func (r *Status) SetReason(reason string) *Status {
	r.Reason = reason
	return r
}

// String prints status info.
func (r *Status) String() string {
	if r == nil {
		return "<nil>"
	}
	b, _ := r.MarshalJSON()
	return goutil.BytesToString(b)
}

// MarshalJSON marshals Status into JSON, implements json.Marshaler interface.
func (r *Status) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	var b = append(reA, strconv.FormatInt(int64(r.Code), 10)...)
	if len(r.Message) > 0 {
		b = append(b, reB...)
		b = append(b, goutil.StringMarshalJSON(r.Message, false)...)
	}
	if len(r.Reason) > 0 {
		b = append(b, reC...)
		b = append(b, goutil.StringMarshalJSON(r.Reason, false)...)
	}
	b = append(b, '}')
	return b, nil
}

// UnmarshalJSON unmarshals a JSON description of self.
func (r *Status) UnmarshalJSON(b []byte) error {
	if r == nil {
		return nil
	}
	s := goutil.BytesToString(b)
	r.Code = int32(gjson.Get(s, "code").Int())
	r.Message = gjson.Get(s, "message").String()
	r.Reason = gjson.Get(s, "reason").String()
	return nil
}

// ToError converts to error interface.
func (r *Status) ToError() error {
	if r == nil || r.Code == OK {
		return nil
	}
	return (*serror)(unsafe.Pointer(r))
}

type serror Status

func (r *serror) Error() string {
	b, _ := r.toStatus().MarshalJSON()
	return goutil.BytesToString(b)
}

func (r *serror) toStatus() *Status {
	return (*Status)(unsafe.Pointer(r))
}
