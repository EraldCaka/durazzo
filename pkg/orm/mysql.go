package orm

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func initMySQL(dsn string) (*sql.DB, error) {
	return sql.Open("mysql", dsn)
}
