package goutil

import (
	"testing"
)

func TestMd5(t *testing.T) {
	b := []byte("1234567890abcdef")
	t.Logf("text: %s, md5: %s", b, Md5(b))
}

func TestEncrypt(t *testing.T) {
	cipherkey := []byte("1234567890abcdef")
	ciphertext := AESEncrypt(cipherkey, []byte("text1234"))
	t.Logf("ciphertext: %s", ciphertext)
	rawtext, err := AESDecrypt(cipherkey, ciphertext)
	t.Logf("rawtext: %s, error: %v", rawtext, err)
	if string(rawtext) != "text1234" {
		t.Fatalf("expect: %s, but get: %s", "text1234", rawtext)
	}
}
