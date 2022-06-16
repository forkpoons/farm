package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository struct {
	Db *sqlx.DB
}

func New(database, conn string) (*Repository, error) {
	db, err := sqlx.Connect(database, conn)
	if err != nil {
		return nil, err
	}
	return &Repository{Db: db}, nil
}

func (r *Repository) ReadTemperature(name string) (float64, error) {
	var temp []float64
	err := r.Db.Select(&temp, "SELECT temp FROM temperatures")
	return temp[len(temp)-1], err
}

func (r *Repository) WriteTemperature(name string, temp float64) error {
	_, err := r.Db.Exec(
		"INSERT INTO temperatures (`name`, `date`, `temp`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		temp,
	)
	return err
}

func (r *Repository) WriteAction(name string, action bool) error {
	_, err := r.Db.Exec(
		"INSERT INTO actions (`name`, `date`, `action`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		action,
	)
	return err
}
