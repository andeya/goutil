package status

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	var err error
	err = WrapError(err, Map{"author": "andeya"})
	assert.Equal(t, nil, err)

	err = errors.New("error text1")
	err = WrapError(err, Map{"author": "andeya"})
	assert.EqualError(t, err, `author=andeya, error=error text1`)

	err = errors.New("error text2")
	err = WrapError(err, Map{"author": []string{"andeya"}})
	assert.EqualError(t, err, `author=[andeya], error=error text2`)

	err = errors.New("error text3")
	err = WrapError(err, Map{"author": struct {
		Name string
		X    int
	}{Name: "andeya"}})
	assert.EqualError(t, err, `author={Name:andeya X:0}, error=error text3`)

	err = WrapError("error text4", Map{"author": "andeya"})
	assert.EqualError(t, err, `author=andeya, error=error text4`)
}
