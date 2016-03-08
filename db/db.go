package db

import (
	"errors"
	"github.com/ManikDV/storage/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	Session *mgo.Session

	Mongo *mgo.DialInfo

	log = utils.SetUpLogger("db")
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
	return false, err
}

type Database struct {
	Id          string `db:"_id"`
	Partitioned bool   `db:"partitioned"`
	Primary     string `db:"primary"`
}

type Shard struct {
	Id   string `shard:"_id"`
	Host string `shard:"host"`
	Tags string `shard:"tags"`
}

type Mongos struct {
	Id      string `mongos:"_id"`
	Ping    string `mongos:"ping"`
	Up      int    `mongos:"up"`
	Waiting bool   `mongos:"waiting"`
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

// Gets information about cluster topology and it's members: mongoses, shards, DBs
func GetClusterStats() (ClusterStats, error) {

	var databases []Database
	var shards []Shard
	var mongoses []Mongos
	var collections []Collection

	configExists, err := DbExists("config")
	if err != nil {
		log.Error(err)
		return ClusterStats{}, err
	}
	if !configExists {
		return ClusterStats{}, errors.New("Config databse doesn't exist, check whether you are connecting to mongos")
	}

	mainDbExists, err := DbExists(utils.DBName)
	if err != nil {
		log.Error(err)
		return ClusterStats{}, err
	}
	if !mainDbExists {
		return ClusterStats{}, errors.New("Main db with name " + utils.DBName + " doesn't exist")
	}
	var session = Session.Clone()
	defer session.Close()

	var configDB = session.DB("config")
	var mainDB = session.DB(utils.DBName)

	// find all databases in cluster
	if err := configDB.C("databases").Find(nil).All(&databases); err != nil {
		log.Error(err)
		return ClusterStats{}, err
	}

	// find all shards in cluster
	if err := configDB.C("shards").Find(nil).All(&shards); err != nil {
		log.Error(err)
		return ClusterStats{}, err
	}

	// find all mongos in cluster
	if err := configDB.C("mongos").Find(nil).All(&mongoses); err != nil {
		log.Error(err)
		return ClusterStats{}, err
	}

	// find all sharded collections in cluster
	colNames, err := mainDB.CollectionNames()
	if err != nil {
		log.Error(err)
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

	return ClusterStats{}, errors.New("Config db doesn't exist")
}

type DbStats struct {
	Raw           bson.M "raw"
	Objects       int    "objects"
	AvgObjectSize int    "avgObjSize"
	DataSize      int    "dataSize"
	StorageSize   int    "storageSize"
	Indexes       int    "indexes"
}

// Gets detailed statistics of certain database
func GetDbStats(dbName string) (DbStats, error) {
	session := Session.Clone()
	db := session.DB(dbName)

	result := DbStats{}
	err := db.Run(bson.D{{"dbStats", 1}, {"scale", 1}}, &result)
	if err != nil {
		return DbStats{}, err
	}

	return result, nil
}
