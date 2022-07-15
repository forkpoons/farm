package greenhouse

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type repo struct {
	ctx context.Context
	db  *sqlx.DB
}

func newDB(ctx context.Context, db *sqlx.DB) *repo {
	return &repo{
		db:  db,
		ctx: ctx,
	}
}

func (r *repo) ReadDevices() ([]Device, error) {
	var dev []Device
	err := r.db.Select(&dev, "SELECT * FROM greenhouse")
	return dev, err
}

func (r *repo) ReadTemperature(name string) (float64, error) {

	var temp []float64
	err := r.db.Select(&temp, "SELECT temp FROM temperatures")
	return temp[len(temp)-1], err
}

func (r *repo) WriteTemperature(name string, temp float64) error {
	_, err := r.db.Exec(
		"INSERT INTO temperatures (`name`, `date`, `temp`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		temp,
	)
	return err
}

func (r *repo) WriteAction(name string, action bool) error {
	_, err := r.db.Exec(
		"INSERT INTO actions (`name`, `date`, `action`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		action,
	)
	return err
}
