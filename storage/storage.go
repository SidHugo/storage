package main

import (
	"net/http"

	"github.com/ManikDV/storage/api"
	"github.com/ManikDV/storage/utils"
	"gopkg.in/mgo.v2"
)

func main() {
	// Setup DB connection
	session, err := mgo.Dial(utils.DBUrl)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	router := api.NewRouter()
	http.ListenAndServe(":8080", router)
}
