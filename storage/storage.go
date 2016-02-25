package main

import (
	"net/http"

	"database/sql"
	"github.com/ManikDV/storage/api"
	"github.com/ManikDV/storage/db"
)

type DbServer struct {
	db *sql.DB
}

func init() {
	db.InitDb()
}

func main() {
	router := api.NewRouter()
	http.ListenAndServe(":8080", router)
}
