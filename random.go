package goutil

import (
	"crypto/rand"
	"encoding/base64"
	mrand "math/rand"
)

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

var (
	encoding   = base64.URLEncoding
	encoder    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
	encoderLen = len(encoder)
	ignoreMap  = map[byte]struct{}{}
)

// SetRandomSeed sets a new padded Encoding defined by the given alphabet.
func SetRandomSeed(encoderSeed string, ignore ...byte) {
	encoding = base64.NewEncoding(encoderSeed)
	ignoreMap = map[byte]struct{}{}
	if len(ignore) > 0 {
		bMap := map[byte]struct{}{}
		for i := 0; i < len(encoderSeed); i++ {
			bMap[encoderSeed[i]] = struct{}{}
		}
		for _, b := range ignore {
			ignoreMap[b] = struct{}{}
			delete(bMap, b)
		}
		encoder = encoder[:0]
		for b := range bMap {
			encoder = append(encoder, b)
		}
	} else {
		encoder = []byte(encoderSeed)
	}
	encoderLen = len(encoder)
}

// RandomString returns a URL-safe, base64 encoded securely generated
// random string. It will panic if the system's secure random number generator
// fails to function correctly.
// The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.
func RandomString(n int) string {
	d := encoding.DecodedLen(n)
	buf := make([]byte, encoding.EncodedLen(d), n)
	encoding.Encode(buf, RandomBytes(d))
	if len(ignoreMap) > 0 {
		var ok bool
		for i, b := range buf {
			if _, ok = ignoreMap[b]; ok {
				buf[i] = encoder[mrand.Intn(encoderLen)]
			}
		}
	}
	for i := n - len(buf); i > 0; i-- {
		buf = append(buf, encoder[mrand.Intn(encoderLen)])
	}
	return BytesToString(buf)
}
