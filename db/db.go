package db

import (
	"github.com/ManikDV/storage/utils"
	"gopkg.in/mgo.v2"
)

func InitDb() (*mgo.Session, error) {
	// Setup DB connection
	session, err := mgo.Dial(utils.DBUrl)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	return session, nil
}
