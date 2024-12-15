package postgres

import (
	"github.com/jmoiron/sqlx"

	// DB adapter
	_ "github.com/lib/pq"
)

// New creates new database connection to a postgres database
func New(psn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", psn)
	if err != nil {

		return nil, err
	}

	if db.Ping() != nil {

		return nil, err
	}

	return db, nil
}
