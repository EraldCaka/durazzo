package orm

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func initPostgres(dsn string) (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}
