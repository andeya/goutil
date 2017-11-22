package goutil

import (
	"crypto/rand"
	"encoding/base64"
	mrand "math/rand"
)

// NewRandom creates a new padded Encoding defined by the given alphabet.
func NewRandom(encoderSeed string, ignore ...byte) *Random {
	r := new(Random)
	r.encoding = base64.NewEncoding(encoderSeed)
	r.ignoreMap = map[byte]struct{}{}
	if len(ignore) > 0 {
		bMap := map[byte]struct{}{}
		for i := 0; i < len(encoderSeed); i++ {
			bMap[encoderSeed[i]] = struct{}{}
		}
		for _, b := range ignore {
			r.ignoreMap[b] = struct{}{}
			delete(bMap, b)
		}
		r.encoder = r.encoder[:0]
		for b := range bMap {
			r.encoder = append(r.encoder, b)
		}
	} else {
		r.encoder = []byte(encoderSeed)
	}
	r.encoderLen = len(r.encoder)
	return r
}

// Random random string creater.
type Random struct {
	encoding   *base64.Encoding
	encoder    []byte
	encoderLen int
	ignoreMap  map[byte]struct{}
}

// RandomString returns a base64 encoded securely generated
// random string. It will panic if the system's secure random number generator
// fails to function correctly.
// The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.
func (r *Random) RandomString(n int) string {
	d := r.encoding.DecodedLen(n)
	buf := make([]byte, r.encoding.EncodedLen(d), n)
	r.encoding.Encode(buf, RandomBytes(d))
	if len(r.ignoreMap) > 0 {
		var ok bool
		for i, b := range buf {
			if _, ok = r.ignoreMap[b]; ok {
				buf[i] = r.encoder[mrand.Intn(r.encoderLen)]
			}
		}
	}
	for i := n - len(buf); i > 0; i-- {
		buf = append(buf, r.encoder[mrand.Intn(r.encoderLen)])
	}
	return BytesToString(buf)
}

const urlEncoder = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var urlRandom = &Random{
	encoding:   base64.URLEncoding,
	encoder:    []byte(urlEncoder),
	encoderLen: len(urlEncoder),
	ignoreMap:  map[byte]struct{}{},
}

// URLRandomString returns a URL-safe, base64 encoded securely generated
// random string. It will panic if the system's secure random number generator
// fails to function correctly.
// The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.
func URLRandomString(n int) string {
	return urlRandom.RandomString(n)
}

// RandomBytes returns securely generated random bytes. It will panic
// if the system's secure random number generator fails to function correctly.
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}
	return b
}
