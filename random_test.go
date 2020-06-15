package goutil

import (
	"crypto/rand"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const tokenLength = 32

// shortReader provides a broken implementation of io.Reader for testing.
type shortReader struct{}

func (sr shortReader) Read(p []byte) (int, error) {
	return len(p) % 2, io.ErrUnexpectedEOF
}

// TestRandomBytes tests the (extremely rare) case that crypto/rand does
// not return the expected number of bytes.
func TestRandomBytes(t *testing.T) {
	// Pioneered by https://github.com/justinas/nosurf
	original := rand.Reader
	rand.Reader = shortReader{}
	defer func() {
		rand.Reader = original
	}()

	var b = make([]byte, tokenLength)
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("RandomBytes did not report a short read: only read %d bytes", len(b))
		}
	}()

	b = RandomBytes(tokenLength)
}

func TestRandomString(t *testing.T) {
	r := NewRandom("0123456789")
	m := map[string]bool{}
	var lock sync.Mutex
	var group sync.WaitGroup
	count := 10000
	group.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			id := r.RandomString(10)
			lock.Lock()
			m[id] = true
			lock.Unlock()
			group.Done()
		}()
	}
	group.Wait()
	if len(m) != count {
		t.Fail()
	}
	var i int
	t.Log("print the top ten...")
	for id := range m {
		i++
		if i > 50 {
			break
		}
		t.Log(id)
	}
}

func TestRandomStringWithTime(t *testing.T) {
	r := URLRandom()
	tm, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", "2038-12-24 13:45:35 +0800 CST")
	unixTs0 := tm.Unix()
	s0, err0 := r.RandomStringWithTime(32, unixTs0)
	assert.NoError(t, err0)
	t.Logf("RandomStringWithTime: %s", s0)
	for i := 0; i < 100; i++ {
		s, err := r.RandomStringWithTime(32, unixTs0)
		assert.NoError(t, err)
		t.Logf("RandomStringWithTime(%d): %s", i, s)
		unixTs, err := r.ParseTime(s)
		assert.NoError(t, err)
		assert.Equal(t, unixTs0, unixTs)
	}
}
