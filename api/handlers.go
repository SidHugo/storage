package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ManikDV/storage/db"
	"github.com/ManikDV/storage/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var log = utils.SetUpLogger("api")

func Ping(w http.ResponseWriter, r *http.Request) {
	log.Info("Ping")

	responseMessage := "Pong!"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		log.Error(err)
	}
}

func CreateSign(w http.ResponseWriter, r *http.Request) {
	log.Info("CreateSign")

	var sign Sign

	if !CheckCredentials(w, r) {
		return
	}

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
	elapsed := time.Since(start).Nanoseconds() / 1000000
	db.AvgWriteQueryTime = (db.AvgWriteQueryTime + elapsed) / 2
	db.LastWriteQueryTime = elapsed
	if elapsed > db.MaxWriteQueryTime {
		db.MaxWriteQueryTime = elapsed
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sign); err != nil {
		log.Error(err)
	}
}

func GetSign(w http.ResponseWriter, r *http.Request) {
	log.Info("GetSign")

	var sign Sign

	if !CheckCredentials(w, r) {
		return
	}

	signName := mux.Vars(r)["signName"]

	session := db.Session.Clone()
	defer session.Close()

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Find(bson.M{"signname": signName}).One(&sign); err != nil {
		log.Error(err)
	}
	elapsed := time.Since(start).Nanoseconds() / 1000000
	db.AvgReadQueryTime = (db.AvgReadQueryTime + elapsed) / 2
	db.LastReadQueryTime = elapsed
	if elapsed > db.MaxReadQueryTime {
		db.MaxReadQueryTime = elapsed
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sign); err != nil {
		log.Error(err)
	}
}

func GetSigns(w http.ResponseWriter, r *http.Request) {
	log.Info("GetSigns")

	var signs Signs

	if !CheckCredentials(w, r) {
		return
	}

	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Find(nil).All(&signs); err != nil {
		log.Error(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(signs); err != nil {
		log.Error(err)
	}
}

func DeleteSign(w http.ResponseWriter, r *http.Request) {
	log.Info("DeleteSigns")

	if !CheckCredentials(w, r) {
		return
	}

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

	if !CheckCredentials(w, r) {
		return
	}

	var user User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Error(err)
		}
	}

	// Create user in DB
	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBUsersCollectionName)
	if err := collection.Insert(&user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
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

// Tests DB speed - inserts values which qty is specified as GET parameter, immediately retrieves them and measures time
func TestDbSpeed(w http.ResponseWriter, r *http.Request) {
	log.Info("TestDbSpeed")
	var result TestStruct

	// parse arguments
	qty, err := strconv.Atoi(mux.Vars(r)["quantity"])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// check for trying to overload us
	if qty > 10000 {
		log.Error(fmt.Sprintf("Specified incorrect number of records: %d", qty))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please, specify correct value under 10000"))
		return
	}

	start := time.Now()

	session := db.Session.Clone()
	defer session.Close()

	database := session.DB("test")
	collection := database.C("testspeed")

	for i := 0; i < qty; i++ {
		subj := TestStruct{Key: string(i), Value: string(i)}

		// first, insert value
		insertStart := time.Now()
		if err := collection.Insert(&subj); err != nil {
			log.Error("Error adding value to collection", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		elapsed := time.Since(insertStart).Nanoseconds() / 1000000
		db.AvgWriteQueryTime = (db.AvgWriteQueryTime + elapsed) / 2
		db.LastWriteQueryTime = elapsed
		if elapsed > db.MaxWriteQueryTime {
			db.MaxWriteQueryTime = elapsed
		}

		// second, retreive it
		findStart := time.Now()
		if err := collection.Find(bson.M{"key": string(i)}).One(&result); err != nil {
			log.Error("Error retreiving value from collection", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		elapsed = time.Since(findStart).Nanoseconds() / 1000000
		db.AvgReadQueryTime = (db.AvgReadQueryTime + elapsed) / 2
		db.LastReadQueryTime = elapsed
		if elapsed > db.MaxReadQueryTime {
			db.MaxReadQueryTime = elapsed
		}

	}
	// clean up
	collection.DropCollection()

	end := time.Since(start)
	log.Infof("Test successfully passed, inserting and retreiving %d values took %d milliseconds", qty, end.Nanoseconds()/1000000)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Test successfully passed, inserting and retreiving %d values took %d milliseconds", qty, end.Nanoseconds()/1000000)))
}

// Processes request for queries stats: average read/write time and last read/write time
func GetQueryStats(w http.ResponseWriter, r *http.Request) {
	log.Info("GetQueryStats")

	queriesStats := QueryStats{db.AvgWriteQueryTime, db.AvgReadQueryTime, db.LastWriteQueryTime, db.LastReadQueryTime, db.MaxWriteQueryTime, db.MaxReadQueryTime}
	log.Info(queriesStats)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(queriesStats); err != nil {
		log.Error(err)
	}
}

// Check header credentials in http request
func CheckCredentials(w http.ResponseWriter, r *http.Request) bool {
	log.Debug("CheckCredentials")

	login := r.Header.Get("login")
	password := r.Header.Get("password")

	if login == "" || password == "" {
		log.Info("Request withour login or password, declining")
		w.WriteHeader(http.StatusForbidden)
		return false
	} else {
		decodedLogin, err := base64.StdEncoding.DecodeString(login)
		decodedPassword, err := base64.StdEncoding.DecodeString(password)
		decryptedLogin, err := utils.AESDecrypt(decodedLogin)
		if err != nil {
			log.Errorf("Decryption failed for sign message: %s, error: %s", login, err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}

		decryptedPassword, err := utils.AESDecrypt(decodedPassword)
		if err != nil {
			log.Errorf("Decryption failed for sign message: %s, error: %s", password, err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
		if string(decryptedLogin[:]) != utils.Conf.AuthLogin || string(decryptedPassword[:]) != utils.Conf.AuthPassword {
			log.Error("Wrong credentials")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Wrong credentials"))
			return false
		}
	}
	return true
}
