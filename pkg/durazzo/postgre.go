package durazzo

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func initPostgres(dsn string) (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}
