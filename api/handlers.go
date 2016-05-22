package api

import (
	"encoding/json"
	"github.com/ManikDV/storage/db"
	"github.com/ManikDV/storage/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
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

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Insert(&sign); err != nil {
		log.Error(err)
	}
	elapsed := time.Since(start)
	db.AvgWriteQueryTime = (db.AvgWriteQueryTime + (elapsed.Nanoseconds() / 1000)) / 2
	db.LastWriteQueryTime = elapsed.Nanoseconds() / 1000

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

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Find(bson.M{"signname": signName}).One(&sign); err != nil {
		log.Error(err)
	}
	elapsed := time.Since(start)
	db.AvgReadQueryTime = (db.AvgReadQueryTime + (elapsed.Nanoseconds() / 1000)) / 2
	db.LastReadQueryTime = elapsed.Nanoseconds() / 1000

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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Info("CreateUser")

	var user User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Error(err)
		}
	}

	// Create user in DB
	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.DBName).C(utils.DBUsersCollectionName)
	if err := collection.Insert(&user); err != nil {
		log.Error(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error(err)
	}
}

// Processes request for cluster info
func GetClusterInfo(w http.ResponseWriter, r *http.Request) {
	log.Info("GetClusterInfo")

	session := db.Session.Clone()
	defer session.Close()

	stats, err := db.GetClusterStats()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Error(err)
		}
	}
}

// Processes request for certain DB info, specified by dbName parameter
func GetDbInfo(w http.ResponseWriter, r *http.Request) {
	log.Info("GetDbInfo")

	dbName := mux.Vars(r)["dbName"]
	session := db.Session.Clone()
	defer session.Close()

	dbStats, err := db.GetDbStats(dbName)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(dbStats); err != nil {
			log.Error(err)
		}
	}
}

// Processes request for queries stats: average read/write time and last read/write time
func GetQueryStats(w http.ResponseWriter, r *http.Request)  {
	log.Info("GetQueryStats")

	queriesStats := QueryStats{db.AvgWriteQueryTime, db.AvgReadQueryTime, db.LastWriteQueryTime, db.LastReadQueryTime}
	log.Info(queriesStats)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(queriesStats); err != nil {
		log.Error(err)
	}
}
// Processes request for cluster info
func GetClusterInfo(w http.ResponseWriter, r *http.Request) {
	log.Info("GetClusterInfo")

	session := db.Session.Clone()
	defer session.Close()

	stats, err := db.GetClusterStats()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Error(err)
		}
	}
}

// Processes request for certain DB info, specified by dbName parameter
func GetDbInfo(w http.ResponseWriter, r *http.Request) {
	log.Info("GetDbInfo")

	dbName := mux.Vars(r)["dbName"]
	session := db.Session.Clone()
	defer session.Close()

	dbStats, err := db.GetDbStats(dbName)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(dbStats); err != nil {
			log.Error(err)
		}
	}
}

// Processes request for queries stats: average read/write time and last read/write time
func GetQueryStats(w http.ResponseWriter, r *http.Request) {
	log.Info("GetQueryStats")

	queriesStats := QueryStats{db.AvgWriteQueryTime, db.AvgReadQueryTime, db.LastWriteQueryTime, db.LastReadQueryTime}
	log.Info(queriesStats)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(queriesStats); err != nil {
		log.Error(err)
	}
}
