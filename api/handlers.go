package api

import (
	"encoding/json"
	"github.com/ManikDV/storage/db"
	"github.com/ManikDV/storage/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	responseMessage := "Pong!"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		panic(err)
	}
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

func GetSigns(w http.ResponseWriter, r *http.Request) {
	var signs Signs

	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.DBName).C(utils.DBCollectionName)
	if err := collection.Find(nil).All(&signs); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusFound)
	if err := json.NewEncoder(w).Encode(signs); err != nil {
		panic(err)
	}
}

func DeleteSign(w http.ResponseWriter, r *http.Request) {
	signName := mux.Vars(r)["signName"]
	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.DBName).C(utils.DBCollectionName)
	if err := collection.Remove(bson.M{"signname": signName}); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(signName); err != nil { // мб тут должно быть не signName, компилятор ругался прост
		panic(err)
	}
}
