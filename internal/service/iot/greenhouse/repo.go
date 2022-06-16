package greenhouse

import (
	"github.com/jmoiron/sqlx"
	"time"
	"context"
)


type repo struct {
	ctx context.Context
	db *sqlx.DB
}

func newDB(ctx context.Context, db *sqlx.DB) *repo {
	return &repo{
		db: db,
		ctx: ctx,
	}
}

func (r *repo) ReadTemperature(name string) (float64, error) {

	var temp []float64
	err := r.db.Select(&temp, "SELECT temp FROM temperatures")
	return temp[len(temp)-1], err
}

func WriteTemperature(name string, temp float64) error {
	_, err := Data.Exec(
		"INSERT INTO temperatures (`name`, `date`, `temp`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		temp,
	)
	return err
}

func WriteAction(name string, action bool) error {
	_, err := Data.Exec(
		"INSERT INTO actions (`name`, `date`, `action`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		action,
	)
	return err
}
