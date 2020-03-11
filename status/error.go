package status

import (
	"fmt"
)

// WrapError wraps an error with fields.
// NOTE:
//  if err==nil, return nil
func WrapError(cause interface{}, fields map[string]interface{}) error {
	err := toErr(cause)
	if err == nil {
		return nil
	}
	return &causeWithFields{
		err:    err,
		fields: fields,
	}
}

type causeWithFields struct {
	err    error
	fields map[string]interface{}
}

func (c *causeWithFields) Error() string {
	var s string
	for k, v := range c.fields {
		s += fmt.Sprintf("%s=%+v, ", k, v)
	}
	s += fmt.Sprintf("error=%s", c.err.Error())
	return s
}
