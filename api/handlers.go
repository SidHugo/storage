package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Pong!\n")
}

func DataCreate(w http.ResponseWriter, r *http.Request) {
	var sign Sign

	body, err := ioutil.ReadAll(io.Reader(r.Body))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &sign); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// Create sign in DB
}
