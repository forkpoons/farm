package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
