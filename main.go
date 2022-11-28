package main

import (
	"database/sql"
	"log"

	"github.com/blessedmadukoma/go-simple-bank/api"
	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable"
	serverAddress = "0.0.0.0:8000"
)

func main() {
	// connect to database
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.StartServer(serverAddress)
	if err != nil {
		log.Fatal("cannot start server!")
	}
}
