package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func New(database, conn string) (*Repository, error) {
	db, err := sqlx.Connect(database, conn)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) WriteTemperature(temp int) error {
	_, err := r.db.Exec(
		"INSERT INTO temperatures (`date`, `temp`) VALUES(?,?)",
		time.Now().Add(time.Hour*7),
		temp,
	)
	return err
}
