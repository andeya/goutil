package goutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetToStrings(t *testing.T) {
	a := make([]string, 0, 10)
	a = append(a, "1", "2", "3")
	b := SetToStrings(a, "4", "2", "5")
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, b)
	assert.Equal(t, []string{"1", "2", "3"}, a)
}

func TestIntersectStrings(t *testing.T) {
	cases := []struct {
		set1, set2, intersect []string
	}{
		{
			[]string{"a", "b", "c", "d"},
			[]string{"b", "d"},
			[]string{"b", "d"},
		},
		{
			[]string{"a", "b", "c", "d", "e", "d"},
			[]string{"b", "d", "d", "d"},
			[]string{"b", "d"},
		},
		{
			[]string{"a"},
			nil,
			[]string{},
		},
		{
			[]string{"a"},
			[]string{"b"},
			[]string{},
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.intersect, IntersectStrings(c.set1, c.set2))
	}
}

func TestStringsToXXX(t *testing.T) {
	var cases = []struct {
		strs     []string
		expected []int
	}{
		{[]string{"1", "", "3"}, []int{1, 0, 3}},
	}
	for _, c := range cases {
		r, err := StringsToInts(c.strs)
		assert.EqualError(t, err, "strconv.Atoi: parsing \"\": invalid syntax")
		r, err = StringsToInts(c.strs, true)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, r)
	}
}
