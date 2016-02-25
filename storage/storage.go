package main

import (
	"net/http"

	"github.com/ManikDV/storage/api"
	"github.com/ManikDV/storage/db"
)

func main() {
	session, err := db.InitDb()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	router := api.NewRouter()
	http.ListenAndServe(":8080", router)
}
