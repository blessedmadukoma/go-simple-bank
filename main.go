package main

import (
	"database/sql"
	"log"

	"github.com/blessedmadukoma/go-simple-bank/api"
	db "github.com/blessedmadukoma/go-simple-bank/db/sqlc"
	"github.com/blessedmadukoma/go-simple-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	// connect to database
	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.StartServer(config.SERVER_ADDRESS)
	if err != nil {
		log.Fatal("cannot start server!")
	}
}
