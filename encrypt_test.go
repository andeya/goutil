package goutil

import (
	"encoding/hex"
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
	plaintext, err := AESDecrypt(cipherkey, ciphertext)
	t.Logf("plaintext: %s, error: %v", plaintext, err)
	if string(plaintext) != "text1234" {
		t.Fatalf("expect: %s, but get: %s", "text1234", plaintext)
	}
}

func TestPading(t *testing.T) {
	plainText := []byte("text1234")
	blockSize := 16
	r := pkcs5Padding(plainText, blockSize)
	dst := make([]byte, hex.EncodedLen(len(r)))
	hex.Encode(dst, r)
	t.Log(string(dst))
}
