package goutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGoTest(t *testing.T) {
	assert.True(t, IsGoTest())
}

func BenchmarkIsGoTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assert.True(b, IsGoTest())
	}
}
