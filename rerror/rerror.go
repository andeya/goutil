// Package rerror is an abstract object of erroneous result.
package rerror

import (
	"encoding/json"
	"strconv"
	"unsafe"

	"github.com/henrylee2cn/goutil"
	"github.com/tidwall/gjson"
)

const (
	// CodeUnknown unknown error code
	CodeUnknown = -1
	// CodeNoError no error code
	CodeNoError = 0 // nil error (ok)
)

// Rerror an abstract object of erroneous result
type Rerror struct {
	// Code error code
	Code int32 `json:"code"`
	// Message the error message displayed to the user (optional)
	Message string `json:"message"`
	// Reason the cause of the error for debugging (optional)
	Reason string `json:"reason"`
}

var (
	_ json.Marshaler   = new(Rerror)
	_ json.Unmarshaler = new(Rerror)

	reA = []byte(`{"code":`)
	reB = []byte(`,"message":`)
	reC = []byte(`,"reason":`)

	rerrUnknown = New(CodeUnknown, "Unknown Error", "")
)

// New creates a *Rerror.
func New(code int32, message, reason string) *Rerror {
	return &Rerror{
		Code:    code,
		Message: message,
		Reason:  reason,
	}
}

// To converts error to *Rerror
func To(err error) *Rerror {
	if err == nil {
		return nil
	}
	r, ok := err.(*rerror)
	if ok {
		return r.toRerror()
	}
	rerr := rerrUnknown.Copy().SetReason(err.Error())
	return rerr
}

// HasError returns true if there are no error.
func (r *Rerror) HasError() bool {
	return r != nil && r.Code != CodeNoError
}

// Copy returns the copy of Rerror
func (r Rerror) Copy() *Rerror {
	return &r
}

// SetMessage sets the error message displayed to the user.
func (r *Rerror) SetMessage(message string) *Rerror {
	r.Message = message
	return r
}

// SetReason sets the cause of the error for debugging.
func (r *Rerror) SetReason(reason string) *Rerror {
	r.Reason = reason
	return r
}

// String prints error info.
func (r *Rerror) String() string {
	if r == nil {
		return "<nil>"
	}
	b, _ := r.MarshalJSON()
	return goutil.BytesToString(b)
}

// MarshalJSON marshals Rerror into JSON, implements json.Marshaler interface.
func (r *Rerror) MarshalJSON() ([]byte, error) {
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
func (r *Rerror) UnmarshalJSON(b []byte) error {
	if r == nil {
		return nil
	}
	s := goutil.BytesToString(b)
	r.Code = int32(gjson.Get(s, "code").Int())
	r.Message = gjson.Get(s, "message").String()
	r.Reason = gjson.Get(s, "reason").String()
	return nil
}

// ToError converts to error
func (r *Rerror) ToError() error {
	if r == nil || r.Code == CodeNoError {
		return nil
	}
	return (*rerror)(unsafe.Pointer(r))
}

type rerror Rerror

func (r *rerror) Error() string {
	b, _ := r.toRerror().MarshalJSON()
	return goutil.BytesToString(b)
}

func (r *rerror) toRerror() *Rerror {
	return (*Rerror)(unsafe.Pointer(r))
}
