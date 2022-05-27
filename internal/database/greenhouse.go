package database

import "time"

func (r *Repository) ReadTemperature(name string) (float64, error) {
	var temp []float64
	err := r.db.Select(&temp, "SELECT temp FROM temperatures")
	return temp[len(temp)-1], err
}

func (r *Repository) WriteTemperature(name string, temp float64) error {
	_, err := r.db.Exec(
		"INSERT INTO temperatures (`name`, `date`, `temp`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		temp,
	)
	return err
}

func (r *Repository) WriteAction(name string, action bool) error {
	_, err := r.db.Exec(
		"INSERT INTO actions (`name`, `date`, `action`) VALUES(?,?,?)",
		name,
		time.Now().Add(time.Hour*7),
		action,
	)
	return err
}
