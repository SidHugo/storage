package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/base64"
	"encoding/json"
	"bytes"
	math "math/rand"
	"time"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1 << letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)
var src = math.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

type Sign2 struct {
	Link         string `json:"link"`
	Base64string string `json:"base64string"`
}

type Signs2 []Sign2

func main() {
	encryptedLogin,err := AESEncrypt([]byte("login"))
	if err != nil {
		fmt.Println(err)
	}
	encryptedPassword,err := AESEncrypt([]byte("password"))
	if err != nil {
		fmt.Println(err)
	}

	var sign Sign2
	CreateSign(encryptedLogin, encryptedPassword, &sign)
	GetSign(encryptedLogin, encryptedPassword, &sign)
	GetSignJson(encryptedLogin, encryptedPassword, &sign)
	GetAllSigns(encryptedLogin, encryptedPassword)
}

func CreateSign(encryptedLogin []byte, encryptedPassword []byte, newSign *Sign2) {
	sign := Sign2{Link: RandStringBytesMaskImprSrc(5), Base64string: RandStringBytesMaskImprSrc(5)}
	signJson, err := json.Marshal(&sign)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/signs", bytes.NewBuffer(signJson))
	req.Header.Set("login", base64.StdEncoding.EncodeToString(encryptedLogin))
	req.Header.Set("password", base64.StdEncoding.EncodeToString(encryptedPassword))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Create sign:")
	fmt.Println(resp.StatusCode)
	fmt.Println(string(contents[:]))
	*newSign = sign
}

func GetAllSigns(encryptedLogin []byte, encryptedPassword []byte) {
	var signs Signs2
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/signs", nil)
	req.Header.Set("login", base64.StdEncoding.EncodeToString(encryptedLogin))
	req.Header.Set("password", base64.StdEncoding.EncodeToString(encryptedPassword))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(contents, &signs); err != nil {
		fmt.Println(err)
	}

	fmt.Println("----------------")
	fmt.Println("Get all signs:")
	fmt.Println(resp.StatusCode)
	fmt.Println(string(contents[:]))
	fmt.Println(signs)
}

func GetSign(encryptedLogin []byte, encryptedPassword []byte, sign *Sign2) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/signs/" + sign.Link, nil)
	req.Header.Set("login", base64.StdEncoding.EncodeToString(encryptedLogin))
	req.Header.Set("password", base64.StdEncoding.EncodeToString(encryptedPassword))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(contents, &sign); err != nil {
		fmt.Println(err)
	}

	fmt.Println("----------")
	fmt.Println("Get sign:")
	fmt.Println(resp.StatusCode)
	fmt.Println(string(contents[:]))
	fmt.Println(sign)
}

func GetSignJson(encryptedLogin []byte, encryptedPassword []byte, sign *Sign2) {
	signJson, err := json.Marshal(&sign)
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/signsJson", bytes.NewBuffer(signJson))
	req.Header.Set("login", base64.StdEncoding.EncodeToString(encryptedLogin))
	req.Header.Set("password", base64.StdEncoding.EncodeToString(encryptedPassword))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var signResp Sign2
	if err := json.Unmarshal(contents, &signResp); err != nil {
		fmt.Println(err)
	}
	fmt.Println("---------------")
	fmt.Println("Get sign JSON:")
	fmt.Println("Status code: " + string(resp.StatusCode))
	fmt.Println(string(contents[:]))
	fmt.Println(signResp)
}

func AESEncrypt(text []byte) (encrypted []byte, err error) {
	// creating new cipher block
	block, err := aes.NewCipher([]byte("abcabcabcaabcabc"))
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
	block, err := aes.NewCipher([]byte("abcabcabcaabcabc"))
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