package status

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	var err error
	err = WrapError(err, Map{"author": "henrylee2cn"})
	assert.Equal(t, nil, err)

	err = errors.New("error text1")
	err = WrapError(err, Map{"author": "henrylee2cn"})
	assert.EqualError(t, err, `author=henrylee2cn, error=error text1`)

	err = errors.New("error text2")
	err = WrapError(err, Map{"author": []string{"henrylee2cn"}})
	assert.EqualError(t, err, `author=[henrylee2cn], error=error text2`)

	err = errors.New("error text3")
	err = WrapError(err, Map{"author": struct {
		Name string
		X    int
	}{Name: "henrylee2cn"}})
	assert.EqualError(t, err, `author={Name:henrylee2cn X:0}, error=error text3`)

	err = WrapError("error text4", Map{"author": "henrylee2cn"})
	assert.EqualError(t, err, `author=henrylee2cn, error=error text4`)
}
