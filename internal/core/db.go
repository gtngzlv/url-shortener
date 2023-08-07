package core

import (
	"database/sql"
	"time"
)

func InitDB(connString, resultURL string) (*sql.DB, string) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, ""
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)
	return db, resultURL
}
