package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var connection *sql.DB

func OpenConnection(path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("error opening connection - %v", err)
	}
	connection = db
}

func CloseConnection() {
	err := connection.Close()
	if err != nil {
		log.Fatalf("error closing connection - %v", err)
	}
}
