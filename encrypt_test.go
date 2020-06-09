package goutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMd5(t *testing.T) {
	b := []byte("1234567890abcdef")
	t.Logf("text: %s, md5: %s", b, Md5(b))
}

func TestSha1(t *testing.T) {
	b := []byte("1234567890abcdef")
	t.Logf("text: %s, sha1: %s", b, Sha1(b))
}

func TestSha256(t *testing.T) {
	b := []byte("1234567890abcdef")
	t.Logf("text: %s, Sha256: %s", b, Sha256(b))
}

func TestSha512(t *testing.T) {
	b := []byte("1234567890abcdef")
	t.Logf("text: %s, Sha512: %s", b, Sha512(b))
}

func TestFnv1aToUint(t *testing.T) {
	b := []byte("1234567890abcdef")
	t.Logf("text: %s, Fnv1aToUint64: %d", b, Fnv1aToUint64(b))
	t.Logf("text: %s, Fnv1aToUint32: %d", b, Fnv1aToUint32(b))
}

var (
	_cipherkey = []byte("1234567890abcdef")
	_plaintext = []byte("text1234")
)

func TestAESEncrypt(t *testing.T) {
	ciphertext := AESEncrypt(_cipherkey, _plaintext)
	t.Logf("ciphertext hex: %s", ciphertext)
	r, err := AESDecrypt(_cipherkey, ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, _plaintext, r)

	ciphertext = AESEncrypt(_cipherkey, _plaintext, true)
	t.Logf("ciphertext base64: %s", ciphertext)
	r, err = AESDecrypt(_cipherkey, ciphertext, true)
	assert.NoError(t, err)
	assert.Equal(t, _plaintext, r)
}

func TestAESCBCEncrypt(t *testing.T) {
	ciphertext := AESCBCEncrypt(_cipherkey, _plaintext)
	t.Logf("ciphertext hex: %s", ciphertext)
	r, err := AESCBCDecrypt(_cipherkey, ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, _plaintext, r)

	ciphertext = AESCBCEncrypt(_cipherkey, _plaintext, true)
	t.Logf("ciphertext base64: %s", ciphertext)
	r, err = AESCBCDecrypt(_cipherkey, ciphertext, true)
	assert.NoError(t, err)
	assert.Equal(t, _plaintext, r)
}

func TestAESCTREncrypt(t *testing.T) {
	ciphertext := AESCTREncrypt(_cipherkey, _plaintext)
	t.Logf("ciphertext hex: %s", ciphertext)
	r, err := AESCTRDecrypt(_cipherkey, ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, _plaintext, r)

	ciphertext = AESCTREncrypt(_cipherkey, _plaintext, true)
	t.Logf("ciphertext base64: %s", ciphertext)
	r, err = AESCTRDecrypt(_cipherkey, ciphertext, true)
	assert.NoError(t, err)
	assert.Equal(t, _plaintext, r)
}

func TestPading(t *testing.T) {
	blockSize := 16
	padded := hexEncode(pkcs5Padding(_plaintext, blockSize))
	t.Log(string(padded))
	r, err := hexDecode(padded)
	assert.NoError(t, err)
	unpaded, err := pkcs5Unpadding(r)
	assert.NoError(t, err)
	assert.Equal(t, _plaintext, unpaded)
}
