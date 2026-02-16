package database

import (
	"github.com/go-pg/pg/v10"
)

func NewConnect() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     "postgres:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "autenticacao",
	})

	return db
}
