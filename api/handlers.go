package api

import (
	"encoding/json"
	"fmt"
	"github.com/ManikDV/storage/db"
	"github.com/ManikDV/storage/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Pong!\n")
}

func CreateSign(w http.ResponseWriter, r *http.Request) {
	var sign Sign

	body, err := ioutil.ReadAll(r.Body)
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
	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.DBName).C(utils.DBCollectionName)
	if err := collection.Insert(&sign); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(sign); err != nil {
		panic(err)
	}
}

func GetSign(w http.ResponseWriter, r *http.Request) {
	var sign Sign

	signName := mux.Vars(r)["signName"]

	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.DBName).C(utils.DBCollectionName)
	if err := collection.Find(bson.M{"signname": signName}).One(&sign); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusFound)
	if err := json.NewEncoder(w).Encode(sign); err != nil {
		panic(err)
	}
}
