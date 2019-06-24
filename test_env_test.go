package goutil

import "testing"

func TestIsGoTest(t *testing.T) {
	if !IsGoTest() {
		t.FailNow()
	}
}

func BenchmarkIsGoTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if !IsGoTest() {
			b.FailNow()
		}
	}
}
