package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DataStore struct {
	Db *sql.DB
}

func NewDataStore(dbName string) (DataStore, error) {
	// get database connection
	Db, err := getConnection(dbName)
	if err != nil {
		return DataStore{}, err
	}

	return DataStore{Db}, nil
}

func getConnection(dbName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)
	// return db if already connected
	if db != nil {
		return db, nil
	}

	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	log.Println("Connected to database")

	return db, nil
}
