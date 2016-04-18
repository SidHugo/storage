package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"crypto/rand"
	"errors"
)

func AESEncrypt(text []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher([]byte(utils.AESKey))
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
	return ciphertext, nil
}

func AESDecrypt(text []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher([]byte(utils.AESKey))
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return text, nil
}

