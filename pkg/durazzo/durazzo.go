package durazzo

import (
	"database/sql"
	logging "github.com/EraldCaka/durazzo/pkg/logs"
	"log/slog"
)

type Durazzo struct {
	Db  *sql.DB
	log *slog.Logger
}

// NewDurazzo creates a Durazzo instance
func NewDurazzo(config Config) *Durazzo {
	db := newConnection(config)

	return &Durazzo{
		Db:  db,
		log: slog.New(logging.NewHandler(nil)).With(slog.Group("db")),
	}
}

func (d *Durazzo) Close() error {
	return d.Db.Close()
}
