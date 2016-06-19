package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ManikDV/storage/api"
	"io/ioutil"
	math "math/rand"
	"net/http"
	"time"
	"github.com/ManikDV/storage/utils"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
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
	encryptedLogin, err := utils.CBCEncrypt([]byte("login"))
	if err != nil {
		fmt.Println(err)
	}
	encryptedPassword, err := utils.CBCEncrypt([]byte("password"))
	if err != nil {
		fmt.Println(err)
	}

	var sign Sign2
	var user api.User
	CreateUser(encryptedLogin, encryptedPassword, &user)
	GetUser(encryptedLogin, encryptedPassword, &user)
	Authorize(encryptedLogin, encryptedPassword, user.Login, user.Password)
	Authorize(encryptedLogin, encryptedPassword, user.Login, "456")
	GetUserIps(encryptedLogin, encryptedPassword, &user)
	CreateSign(encryptedLogin, encryptedPassword, &sign)
	GetSign(encryptedLogin, encryptedPassword, &sign)
	GetSignJson(encryptedLogin, encryptedPassword, &sign)
	GetAllSigns(encryptedLogin, encryptedPassword)
}

func Authorize(encryptedLogin []byte, encryptedPassword []byte, inlogin string, inpass string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:8080/users/authorize/%s/%s", inlogin, inpass), nil)
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
	fmt.Println("------------")
	fmt.Println("Authorize:")
	fmt.Printf("http://127.0.0.1:8080/users/authorize/%s/%s\n", inlogin, inpass)
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Authorization successfull")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization unsuccessfull")
	} else {
		fmt.Printf("Unknown code %d. probably you should check it\n", resp.StatusCode)
	}
	fmt.Println(string(contents[:]))
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
	fmt.Println("-------------")
	fmt.Println("Create sign:")
	fmt.Println(resp.StatusCode)
	fmt.Println(string(contents[:]))
	*newSign = sign
}

func CreateUser(encryptedLogin []byte, encryptedPassword []byte, newUser *api.User) {
	map1 := make(map[string]int)
	map1["1"] = 1
	array1 := []string{"ip1", "ip2", "ip3"}
	user := api.User{Key: math.Intn(1000),
		Login:          RandStringBytesMaskImprSrc(5),
		Password:       RandStringBytesMaskImprSrc(5),
		PreviusResults: array1,
		SubscribersIP:  array1,
		Subscriptions:  map1,
	}
	userJson, err := json.Marshal(&user)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/users", bytes.NewBuffer(userJson))
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
	fmt.Println("-------------")
	fmt.Println("Create user:")
	fmt.Println(resp.StatusCode)
	fmt.Println(string(contents[:]))
	*newUser = user
}

func GetUserIps(encryptedLogin []byte, encryptedPassword []byte, user *api.User) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:8080/users/getips/%s", user.Login), nil)
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

	type Answer struct {
		Results []string `json:"address"`
	}
	var answer Answer
	if err := json.Unmarshal(contents, &answer); err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("----------")
	fmt.Println("Get user ips:")
	fmt.Printf("http://127.0.0.1:8080/users/getips/%s\n", user.Login)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(contents[:]))
	fmt.Println(answer)
}

func GetUser(encryptedLogin []byte, encryptedPassword []byte, user *api.User) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:8080/users/%d", user.Key), nil)
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

	var newUser api.User
	if err := json.Unmarshal(contents, &newUser); err != nil {
		fmt.Println(err)
	}

	fmt.Println("----------")
	fmt.Println("Get certain user:")
	fmt.Printf("http://127.0.0.1:8080/users/%d\n", user.Key)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(contents[:]))
	fmt.Println(newUser)
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
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/signs/"+sign.Link, nil)
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


