package orm

import (
	"database/sql"
	"fmt"
	"log"
)

type Config struct {
	Driver string
	DSN    string
}

func NewConnection(config Config) *sql.DB {
	var db *sql.DB
	var err error

	switch config.Driver {
	case Sqlite:
		db, err = initSQLite(config.DSN)
	case Postgres:
		db, err = initPostgres(config.DSN)
	case Mysql:
		db, err = initMySQL(config.DSN)
	default:
		log.Fatalf("Unsupported driver: %s", config.Driver)
	}

	if err != nil {
		log.Fatalf("Failed to connect to %s database: %v", config.Driver, err)
	}
	fmt.Printf("Connected to %s database successfully!\n", config.Driver)
	return db
}
