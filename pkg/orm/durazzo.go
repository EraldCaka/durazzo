package orm

import (
	"database/sql"
)

type Durazzo struct {
	db         *sql.DB
	conditions []string
	args       []interface{}
	limit      int
}

// NewDurazzo creates a Durazzo instance
func NewDurazzo(config Config) *Durazzo {
	db := NewConnection(config)
	return &Durazzo{
		db: db,
	}
}
