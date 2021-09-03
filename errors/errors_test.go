package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	err1 := errors.New("error_text1")
	err2 := errors.New("error not include")
	err3 := New("error_text_3")
	errs := []error{
		errors.New("error_text2"),
		nil,
		errors.New("error_text4"),
		err3,
		err1,
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

	is := errors.Is(err, err1)
	assert.True(t, is)
	is = errors.Is(err, err2)
	assert.False(t, is)

	target := &myerror{}
	as := errors.As(err, &target)
	assert.True(t, as)
	assert.Equal(t, target, err3)
}
