package api

import (
	"encoding/base64"
	"fmt"
	"github.com/ManikDV/storage/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test correct credentials
func TestCheckCredentialsCorrect(t *testing.T) {
	utils.SetConfig("../config.toml")
	req, err := http.NewRequest("GET", "localhost:8080/createSign", nil)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	encryptedLogin, err := utils.CBCEncrypt([]byte(utils.Conf.AuthLogin))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	encryptedPassword, err := utils.CBCEncrypt([]byte(utils.Conf.AuthPassword))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	req.Header.Set("login", base64.StdEncoding.EncodeToString(encryptedLogin[:]))
	req.Header.Set("password", base64.StdEncoding.EncodeToString(encryptedPassword[:]))

	w := httptest.NewRecorder()
	result := CheckCredentials(w, req)
	if !result {
		t.FailNow()
	}
}

// Test empty credentials
func TestCheckCredentialsEmpty(t *testing.T) {
	utils.SetConfig("../config.toml")
	req, err := http.NewRequest("GET", "localhost:8080/createSign", nil)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	w := httptest.NewRecorder()
	result := CheckCredentials(w, req)
	if result {
		t.FailNow()
	}
}

// Test wrong credentials
func TestCheckCredentialsIncorrect(t *testing.T) {
	utils.SetConfig("../config.toml")
	req, err := http.NewRequest("GET", "localhost:8080/createSign", nil)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	encryptedLogin, err := utils.CBCEncrypt([]byte("test"))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	encryptedPassword, err := utils.CBCEncrypt([]byte("test"))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	req.Header.Set("login", string(encryptedLogin[:]))
	req.Header.Set("password", string(encryptedPassword[:]))

	w := httptest.NewRecorder()
	result := CheckCredentials(w, req)
	if result {
		t.FailNow()
	}
}
