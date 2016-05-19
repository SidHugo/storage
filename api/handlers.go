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

var log = utils.SetUpLogger("api")

func Ping(w http.ResponseWriter, r *http.Request) {
	log.Info("Ping")

	responseMessage := "Pong!"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		log.Error(err)
	}
}

func CreateSign(w http.ResponseWriter, r *http.Request) {
	log.Info("CreateSign")

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
			log.Error(err)
		}
	}

	// Create sign in DB
	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Insert(&sign); err != nil {
		log.Error(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(sign); err != nil {
		log.Error(err)
	}
}

func GetSign(w http.ResponseWriter, r *http.Request) {
	log.Info("GetSign")

	var sign Sign

	signName := mux.Vars(r)["signName"]

	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Find(bson.M{"signname": signName}).One(&sign); err != nil {
		log.Error(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusFound)
	if err := json.NewEncoder(w).Encode(sign); err != nil {
		log.Error(err)
	}
}

func GetSigns(w http.ResponseWriter, r *http.Request) {
	log.Info("GetSigns")

	var signs Signs

	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Find(nil).All(&signs); err != nil {
		log.Error(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusFound)
	if err := json.NewEncoder(w).Encode(signs); err != nil {
		log.Error(err)
	}
}

func DeleteSign(w http.ResponseWriter, r *http.Request) {
	log.Info("DeleteSigns")

	signName := mux.Vars(r)["signName"]
	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Remove(bson.M{"signname": signName}); err != nil {
		log.Error(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(signName); err != nil {
		log.Error(err)
	}
}
