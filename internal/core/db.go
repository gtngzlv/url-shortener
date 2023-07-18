package core

import (
	"database/sql"
)

func InitDB(connString, resultURL string) (*sql.DB, string) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, ""
	}
	return db, resultURL
}
