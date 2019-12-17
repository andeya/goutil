package goutil

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// Md5 returns the MD5 checksum string of the data.
func Md5(b []byte) string {
	checksum := md5.Sum(b)
	return hex.EncodeToString(checksum[:])
}

// AESEncrypt uses ECB mode to encrypt a piece of data.
// The cipherkey argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESEncrypt(cipherkey, plaintext []byte) []byte {
	block, err := aes.NewCipher(cipherkey)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	plaintext = pkcs5Padding(plaintext, blockSize)
	r := make([]byte, len(plaintext))
	dst := r
	for len(plaintext) > 0 {
		block.Encrypt(dst, plaintext)
		plaintext = plaintext[blockSize:]
		dst = dst[blockSize:]
	}
	dst = make([]byte, hex.EncodedLen(len(r)))
	hex.Encode(dst, r)
	return dst
}

// AESDecrypt uses ECB mode to decrypt a piece of data.
// The cipherkey argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESDecrypt(cipherkey, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(cipherkey)
	if err != nil {
		return nil, err
	}
	src := make([]byte, hex.DecodedLen(len(ciphertext)))
	_, err = hex.Decode(src, ciphertext)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	r := make([]byte, len(src))
	dst := r
	for len(src) > 0 {
		block.Decrypt(dst, src)
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	return pkcs5Unpadding(r)
}

func pkcs5Padding(plaintext []byte, blockSize int) []byte {
	n := byte(blockSize - len(plaintext)%blockSize)
	for i := byte(0); i < n; i++ {
		plaintext = append(plaintext, n)
	}
	return plaintext
}

func pkcs5Unpadding(r []byte) ([]byte, error) {
	l := len(r)
	if l == 0 {
		return []byte{}, errors.New("input padded bytes is empty")
	}
	last := int(r[l-1])
	n := byte(last)
	pad := r[l-last : l]
	isPad := true
	for _, v := range pad {
		if v != n {
			isPad = false
			break
		}
	}
	if !isPad {
		return r, errors.New("remove pad error")
	}
	return r[:l-last], nil
}
