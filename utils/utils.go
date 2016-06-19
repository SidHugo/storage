package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AES CFB encryption of byte array
func CFBEncrypt(text []byte) (encrypted []byte, err error) {
	// creating new cipher block
	block, err := aes.NewCipher([]byte(Conf.AESKey))
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

// AES CFB decryption of byte array
func CFBDecrypt(text []byte) (decrypted []byte, err error) {
	// creating new cipher block
	block, err := aes.NewCipher([]byte(Conf.AESKey))
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

// AES CBC decryption of byte array
func CBCDecrypt(ciphertext []byte) ([]byte, error) {
	key := []byte(Conf.AESKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)
	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.
	//fmt.Println("BEFORE CUTTING")
	//fmt.Println(ciphertext)
	ciphertext=unpad(ciphertext)
	return ciphertext, nil
}

// AES CBC encryption of byte array
func CBCEncrypt(plaintext []byte) ([]byte, error){
	key := []byte(Conf.AESKey)
	plaintext=pad(plaintext)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	return ciphertext, nil
}

// Padding byte array (for AES CBC). Pads array with zeros and last byte
// indicating length of padding
func pad(in []byte) []byte {
	//length:=len(in)
	padding := aes.BlockSize - (len(in) % aes.BlockSize)
	if padding == 0 {
		padding = aes.BlockSize
	}
	for i := 0; i < padding-1; i++ {
		in = append(in, byte(0))
	}
	in=append(in, byte(padding))
	return in
}

// Unpading byte array (for AES CBC)
func unpad(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}

	length := in[len(in)-1]
	return in[:len(in)-int(length)]
}