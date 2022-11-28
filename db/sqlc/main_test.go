package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/blessedmadukoma/go-simple-bank/util"
	_ "github.com/lib/pq"
)


var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not load config:", err)
	}

	testDB, err = sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
