package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"crypto/rand"
	"errors"
)

func AESEncrypt(text []byte) (encrypted []byte, err error) {
	// creating new cipher block
	block, err := aes.NewCipher([]byte(AESKey))
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	// initialization vector. Client must have the same IV
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	// we use AES CFB
	cfb := cipher.NewCFBEncrypter(block, iv)
	// XORs each byte in the given slice with a byte from the
	// cipher's key stream
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
	return ciphertext, nil
}

func AESDecrypt(text []byte) (decrypted []byte, err error) {
	// creating new cipher block
	block, err := aes.NewCipher([]byte(AESKey))
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	// initialization vector. Client must have the same IV
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	// we use AES CFB
	cfb := cipher.NewCFBDecrypter(block, iv)
	// XORs each byte in the given slice with a byte from the
	// cipher's key stream
	cfb.XORKeyStream(text, text)
	return text, nil
}

