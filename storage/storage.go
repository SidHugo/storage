package main

import (
	"net/http"

	"database/sql"
	"fmt"
	"github.com/ManikDV/storage/api"
	"github.com/ManikDV/storage/db"
)

type DbServer struct {
	db *sql.DB
}

func init() {
	db.InitDb()
}

func PrintClusterInfo() {
	clusterStats, err := db.GetClusterStats()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Mongoses:")
	for _, mongos := range clusterStats.Mongoses {
		fmt.Printf("\tId:%s, ping: %s, uptime:%d, isWaiting:%t \n", mongos.Id, mongos.Ping, mongos.Up, mongos.Waiting)
	}

	fmt.Println("Shards:")
	for _, shard := range clusterStats.Shards {
		fmt.Printf("\tId:%s, host:%s, tags:%s\n", shard.Id, shard.Host, shard.Tags)
	}

	fmt.Println("Databases:")
	for _, database := range clusterStats.Databases {
		fmt.Printf("\tName:%s, partitioned:%t, primary:%s\n", database.Id, database.Partitioned, database.Primary)
	}

	fmt.Println("Collections:")
	for _, collection := range clusterStats.Collections {
		fmt.Printf("\tCollection name:%s, collection count:%d\n", collection.Name, collection.Count)
	}
}

func main() {

	db.InitDb()

	PrintClusterInfo()

	result, err := db.GetDbStats("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

	router := api.NewRouter()
	http.ListenAndServe(":8080", router)
}
