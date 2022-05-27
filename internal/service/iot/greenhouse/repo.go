package greenhouse

import (
	"github.com/forkpoons/farm/internal/database"
	"time"
)

var data = database.Repository{}

func ReadTemperature(name string) (float64, error) {

	var temp []float64
	err := db.Select(&temp, "SELECT temp FROM temperatures")
	return temp[len(temp)-1], err
}

func WriteTemperature(name string, temp float64) error {
	_, err := r.db.Exec(
		"INSERT INTO temperatures (`name`, `date`, `temp`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		temp,
	)
	return err
}

func WriteAction(name string, action bool) error {
	_, err := r.db.Exec(
		"INSERT INTO actions (`name`, `date`, `action`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		action,
	)
	return err
}
