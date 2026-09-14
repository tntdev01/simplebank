package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"simple_bank/util"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres_user:postgres_password@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../.")
	if err != nil {
		log.Fatal("Unable to load config: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Unable to connect to the database: ", err)
	}

	// Initialize testQueries here, for example:
	testQueries = New(testDB)
	os.Exit(m.Run())
}
