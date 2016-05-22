package main

import (
	"net/http"

	"database/sql"
	"github.com/ManikDV/storage/api"
	"github.com/ManikDV/storage/db"
	"github.com/ManikDV/storage/utils"
)

type DbServer struct {
	db *sql.DB
}

var log = utils.SetUpLogger("storage")

func PrintClusterInfo() {
	clusterStats, err := db.GetClusterStats()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Mongoses:")
	for _, mongos := range clusterStats.Mongoses {
		log.Infof("\tId:%s, ping: %s, uptime:%d, isWaiting:%t \n", mongos.Id, mongos.Ping, mongos.Up, mongos.Waiting)
	}

	log.Info("Shards:")
	for _, shard := range clusterStats.Shards {
		log.Infof("\tId:%s, host:%s, tags:%s\n", shard.Id, shard.Host, shard.Tags)
	}

	log.Info("Databases:")
	for _, database := range clusterStats.Databases {
		log.Infof("\tName:%s, partitioned:%t, primary:%s\n", database.Id, database.Partitioned, database.Primary)
	}

	log.Info("Collections:")
	for _, collection := range clusterStats.Collections {
		log.Infof("\tCollection name:%s, collection count:%d\n", collection.Name, collection.Count)
	}
}

func main() {
	log.Info("Application starting")

	utils.SetConfig()

	db.InitDb()

	PrintClusterInfo()

	result, err := db.GetDbStats("test")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(result)

	router := api.NewRouter()
	http.ListenAndServe(":8080", router)

	log.Info("Application ready to receive requests")
}
