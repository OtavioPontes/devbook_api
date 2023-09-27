package database

import (
	"database/sql"
	"devbook_api/src/config"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.ConectionString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
