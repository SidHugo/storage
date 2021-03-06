package api

import (
	"bytes"
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

// Creates new sign in DB
func CreateSign(w http.ResponseWriter, r *http.Request) {
	log.Info("-> CreateSign")

	var sign Sign2

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

	// estimate size
	stats, err := db.GetDbStats(utils.Conf.DBName)
	if err != nil {
		log.Error("Could not extract db info:", err)
	} else {
		log.Infof("New value will be taking %f part of all storage amount", float32(len(body))/float32(stats.StorageSize))
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

	// check for duplicates
	cnt, err := collection.Find(bson.M{"link": sign.Link}).Count()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if cnt > 0 {
		log.Warningf("Cant insert sign %s - already present", sign.Link)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("Cant insert sign %s - already present", sign.Link)))
		return
	}

	start := time.Now()
	if err := collection.Insert(&sign); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
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

// Gets sign specified by URL parameter
func GetSign(w http.ResponseWriter, r *http.Request) {
	log.Info("-> GetSign")

	var sign Sign2

	if !CheckCredentials(w, r) {
		return
	}

	link := mux.Vars(r)["link"]

	session := db.Session.Clone()
	defer session.Close()

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Find(bson.M{"link": link}).One(&sign); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
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

// Gets sign specified by JSON parameter
func GetSignJson(w http.ResponseWriter, r *http.Request) {
	log.Info("-> GetSignJson")

	var sign Sign2

	if !CheckCredentials(w, r) {
		return
	}

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

	if err := json.Unmarshal(body, &sign); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	session := db.Session.Clone()
	defer session.Close()

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Find(bson.M{"link": sign.Link}).One(&sign); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
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

// Gets all signs from DB
func GetSigns(w http.ResponseWriter, r *http.Request) {
	log.Info("-> GetSigns")

	var signs Signs2

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

// Deletes sign by name, specified by URL parameter
func DeleteSign(w http.ResponseWriter, r *http.Request) {
	log.Info("-> DeleteSigns")

	if !CheckCredentials(w, r) {
		return
	}

	link := mux.Vars(r)["link"]
	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBCollectionName)
	if err := collection.Remove(bson.M{"link": link}); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(link); err != nil {
		log.Error(err)
	}
}

// Creates new user
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
	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBUsersCollectionName)
	if err := collection.Insert(&user); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	elapsed := time.Since(start).Nanoseconds() / 1000000
	db.AvgWriteQueryTime = (db.AvgWriteQueryTime + elapsed) / 2
	db.LastWriteQueryTime = elapsed
	if elapsed > db.MaxWriteQueryTime {
		db.MaxWriteQueryTime = elapsed
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error(err)
	}
}

// Gets user by its key and sends it back in json
func GetUser(w http.ResponseWriter, r *http.Request) {
	log.Info("GetUser")

	var user User

	if !CheckCredentials(w, r) {
		return
	}

	userKey := mux.Vars(r)["key"]

	session := db.Session.Clone()
	defer session.Close()

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBUsersCollectionName)
	res, err := strconv.Atoi(userKey)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := collection.Find(bson.M{"key": res}).One(&user); err != nil {
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
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error(err)
	}
}

// Deletes user by its key
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Info("DeleteUser")

	if !CheckCredentials(w, r) {
		return
	}

	userKey := mux.Vars(r)["key"]
	session := db.Session.Clone()
	defer session.Close()

	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBUsersCollectionName)
	res, err := strconv.Atoi(userKey)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := collection.Remove(bson.M{"key": res}); err != nil {
		log.Error(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userKey); err != nil {
		log.Error(err)
	}
}

// Authorize user and sends back 0 (if failed) or 1 (otherwise) in json
func Authorization(w http.ResponseWriter, r *http.Request) {
	log.Info("Authorization")

	var user User

	if !CheckCredentials(w, r) {
		return
	}

	userLogin := mux.Vars(r)["login"]
	userPassword := mux.Vars(r)["password"]

	session := db.Session.Clone()
	defer session.Close()

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBUsersCollectionName)
	if err := collection.Find(bson.M{"login": userLogin}).One(&user); err != nil {
		log.Error(err)
	}
	elapsed := time.Since(start).Nanoseconds() / 1000000
	db.AvgReadQueryTime = (db.AvgReadQueryTime + elapsed) / 2
	db.LastReadQueryTime = elapsed
	if elapsed > db.MaxReadQueryTime {
		db.MaxReadQueryTime = elapsed
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	type tempAnswer struct {
		Answer int `json:"answer"`
	}
	var answer tempAnswer
	if user.Password == userPassword {
		w.WriteHeader(http.StatusOK)
		answer.Answer = 1
		if err := json.NewEncoder(w).Encode(answer); err != nil {
			log.Error(err)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		answer.Answer = 0
		if err := json.NewEncoder(w).Encode(answer); err != nil {
			log.Error(err)
		}
	}
}

// Gets user subscriptions and sends it back in json
func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	log.Info("GetSubscriptions")

	var user User

	if !CheckCredentials(w, r) {
		return
	}

	userLogin := mux.Vars(r)["login"]
	userPassword := mux.Vars(r)["password"]

	session := db.Session.Clone()
	defer session.Close()

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBUsersCollectionName)
	if err := collection.Find(bson.M{"login": userLogin}).One(&user); err != nil {
		log.Error(err)
	}
	elapsed := time.Since(start).Nanoseconds() / 1000000
	db.AvgReadQueryTime = (db.AvgReadQueryTime + elapsed) / 2
	db.LastReadQueryTime = elapsed
	if elapsed > db.MaxReadQueryTime {
		db.MaxReadQueryTime = elapsed
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	type tempAnswer struct {
		Access int            `json:"access"`
		Subs   map[string]int `json:"subs"`
	}
	var answer tempAnswer

	if user.Password == userPassword {
		w.WriteHeader(http.StatusOK)
		answer.Access = 1
		answer.Subs = user.Subscriptions
		if err := json.NewEncoder(w).Encode(answer); err != nil {
			log.Error(err)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		answer.Access = 0
		answer.Subs = nil
		if err := json.NewEncoder(w).Encode(answer); err != nil {
			log.Error(err)
		}
	}
}

// Gets last results and sends it back in json
func GetLastResults(w http.ResponseWriter, r *http.Request) {
	log.Info("GetLastResults")

	var user User

	if !CheckCredentials(w, r) {
		return
	}

	userLogin := mux.Vars(r)["login"]
	userPassword := mux.Vars(r)["password"]
	requiredLogin := mux.Vars(r)["requiredLogin"]

	session := db.Session.Clone()
	defer session.Close()

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBUsersCollectionName)
	if err := collection.Find(bson.M{"login": userLogin}).One(&user); err != nil {
		log.Error(err)
	}
	elapsed := time.Since(start).Nanoseconds() / 1000000
	db.AvgReadQueryTime = (db.AvgReadQueryTime + elapsed) / 2
	db.LastReadQueryTime = elapsed
	if elapsed > db.MaxReadQueryTime {
		db.MaxReadQueryTime = elapsed
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	type tempAnswer struct {
		Access  int      `json:"access"`
		Results []string `json:"results"`
	}
	var answer tempAnswer
	if user.Password == userPassword && user.Subscriptions[requiredLogin] != 0 {
		start = time.Now()
		if err := collection.Find(bson.M{"login": requiredLogin}).One(&user); err != nil {
			log.Error(err)
			return
		}
		elapsed := time.Since(start).Nanoseconds() / 1000000
		db.AvgReadQueryTime = (db.AvgReadQueryTime + elapsed) / 2
		db.LastReadQueryTime = elapsed
		if elapsed > db.MaxReadQueryTime {
			db.MaxReadQueryTime = elapsed
		}

		w.WriteHeader(http.StatusOK)
		answer.Access = 1
		answer.Results = user.PreviusResults
		if err := json.NewEncoder(w).Encode(answer); err != nil {
			log.Error(err)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		answer.Access = 0
		answer.Results = nil
		if err := json.NewEncoder(w).Encode(answer); err != nil {
			log.Error(err)
		}
	}
}

// Gets IP of subscribers and sends it back in json
func GetUserIPs(w http.ResponseWriter, r *http.Request) {
	log.Info("GetUserIPs")

	var user User

	if !CheckCredentials(w, r) {
		return
	}

	userLogin := mux.Vars(r)["login"]

	session := db.Session.Clone()
	defer session.Close()

	start := time.Now()
	collection := session.DB(utils.Conf.DBName).C(utils.Conf.DBUsersCollectionName)
	if err := collection.Find(bson.M{"login": userLogin}).One(&user); err != nil {
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
	type tempAnswer struct {
		Results []string `json:"address"`
	}
	var answer tempAnswer
	answer.Results = user.SubscribersIP
	if err := json.NewEncoder(w).Encode(answer); err != nil {
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

	session := db.Session.Clone()
	defer session.Close()

	database := session.DB("test")

	var iterations = 10
	var sumTime int64 = 0

	for j := 0; j < iterations; j++ {
		start := time.Now()
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
		// count time
		end := time.Since(start)
		sumTime += end.Nanoseconds() / 1000000
	}
	log.Infof("Test successfully passed, inserting and retreiving %d values in %d series took %d milliseconds in average",
		qty, iterations, sumTime/int64(iterations))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Test successfully passed, inserting and retreiving %d values in %d series took %d milliseconds in average",
		qty, iterations, sumTime/int64(iterations))))
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
		decryptedLogin, err := utils.CBCDecrypt(decodedLogin)
		if err != nil {
			log.Errorf("Decryption failed for sign message: %s, error: %s", login, err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}

		decryptedPassword, err := utils.CBCDecrypt(decodedPassword)
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

// returns list of all methods
func GetHelp(w http.ResponseWriter, r *http.Request) {
	log.Info("-> GetHelp")
	var buffer bytes.Buffer
	buffer.WriteString("Available methods:\n")

	for _, routehelp := range textRoutes {
		buffer.WriteString(routehelp + "\n")
	}

	_, err := w.Write(buffer.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
	} else {
		w.WriteHeader(http.StatusOK)
		log.Info("GetHelp success")
	}
}
