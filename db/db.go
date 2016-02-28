package db

import (
	"errors"
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

func DbExists(name string) (bool, error) {
	// get all names
	var session = Session.Clone()
	defer session.Close()

	names, err := session.DatabaseNames()
	if err != nil {
		return false, err
	}

	for _, dbName := range names {
		if dbName == name {
			return true, nil
		}
	}
	return false, nil
}

type Database struct {
	Id          string "_id"
	Partitioned bool   "partitioned"
	Primary     string "primary"
}

type Shard struct {
	Id   string "_id"
	Host string "host"
	Tags string "tags"
}

type Mongos struct {
	Id      string "_id"
	Ping    string "ping"
	Up      int    "up"
	Waiting bool   "waiting"
}

type Collection struct {
	Name  string
	Count int
}

type ClusterStats struct {
	Databases   []Database
	Shards      []Shard
	Mongoses    []Mongos
	Collections []Collection
}

func GetClusterStats() (ClusterStats, error) {

	var databases []Database
	var shards []Shard
	var mongoses []Mongos
	var collections []Collection

	configExists, err := DbExists("config")
	if err != nil {
		return ClusterStats{}, err
	}

	mainDbExists, err := DbExists(utils.DBName)
	if err != nil {
		return ClusterStats{}, err
	}

	if configExists && mainDbExists {
		var session = Session.Clone()
		defer session.Close()

		var configDB = session.DB("config")
		var mainDB = session.DB(utils.DBName)

		// find all databases in cluster
		if err := configDB.C("databases").Find(nil).All(&databases); err != nil {
			// panic(err)
			return ClusterStats{}, err
		}

		// find all shards in cluster
		if err := configDB.C("shards").Find(nil).All(&shards); err != nil {
			// panic(err)
			return ClusterStats{}, err
		}

		// find all mongos in cluster
		if err := configDB.C("mongos").Find(nil).All(&mongoses); err != nil {
			// panic(err)
			return ClusterStats{}, err
		}

		// find all sharded collections in cluster
		colNames, err := mainDB.CollectionNames()
		if err != nil {
			return ClusterStats{}, err
		}

		for _, colName := range colNames {
			colCount, err := mainDB.C(colName).Count()
			if err != nil {
				return ClusterStats{}, err
			}
			collections = append(collections, Collection{colName, colCount})
		}

		return ClusterStats{databases, shards, mongoses, collections}, nil
	}

	return ClusterStats{}, errors.New("Config db doesn't exist")
}
