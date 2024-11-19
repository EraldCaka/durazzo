package durazzo

import (
	"database/sql"
)

type Durazzo struct {
	Db *sql.DB
}

// NewDurazzo creates a Durazzo instance
func NewDurazzo(config Config) *Durazzo {
	db := newConnection(config)
	return &Durazzo{
		Db: db,
	}
}

func (d *Durazzo) Close() error {
	return d.Db.Close()
}
