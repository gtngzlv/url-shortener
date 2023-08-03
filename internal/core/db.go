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
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Second * 10)
	return db, resultURL
}
