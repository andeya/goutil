package errors

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	errs := []error{
		errors.New("error_text1"),
		errors.New("error_text2"),
		nil,
		errors.New("error_text4"),
		errors.New("error_text5"),
		nil,
		errors.New("error_text7"),
	}
	err := Merge(errs...)
	t.Log(err)
	err = Append(err, nil)
	t.Log(err)
	err = Append(err, errs...)
	t.Log(err)
}
