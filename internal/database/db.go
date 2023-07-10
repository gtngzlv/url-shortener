package database

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PostgresDB struct {
	log zap.SugaredLogger
	db  *sqlx.DB
}

func (p PostgresDB) Save(fullURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDB) Get(shortURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDB) Ping() error {
	if err := p.db.Ping(); err != nil {
		return err
	}
	return nil
}

func Init(log zap.SugaredLogger, path string) *PostgresDB {
	db, err := sqlx.Open("postgres", path)
	if err != nil {
		return nil
	}
	return &PostgresDB{
		log: log,
		db:  db,
	}
}
