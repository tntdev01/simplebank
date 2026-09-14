package main

import (
	"database/sql"
	"log"

	api "simple_bank/api"
	db "simple_bank/db/sqlc"
	utils "simple_bank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Unable to load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Unable to connect to the database: ", err)
	}

	store := db.NewStore(conn)
	defer conn.Close()

	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Unable to start the server: ", err)
	}
}
