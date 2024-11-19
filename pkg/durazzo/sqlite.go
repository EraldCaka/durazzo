package durazzo

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func initSQLite(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite3", dsn)
}
