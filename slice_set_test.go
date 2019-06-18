package goutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			[]string{"b", "d", "d"},
		},
		{
			[]string{"a"},
			nil,
			nil,
		},
		{
			[]string{"a"},
			[]string{"b"},
			nil,
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.intersect, IntersectStrings(c.set1, c.set2))
	}
}
