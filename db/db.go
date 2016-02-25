package db

import (
	"github.com/ManikDV/storage/utils"
	"gopkg.in/mgo.v2"
)

var (
	Session *mgo.Session

	Mongo *mgo.DialInfo
)

func InitDb() {
	// Setup DB connection
	mongo, err := mgo.ParseURL(utils.DBUrl)
	if err != nil {
		panic(err)
	}
	session, err := mgo.Dial(utils.DBUrl)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	Session = session
	Mongo = mongo
}
