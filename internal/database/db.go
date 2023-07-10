package database

import (
	"github.com/gtngzlv/url-shortener/internal/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresDB struct {
	log zap.SugaredLogger
	db  *sqlx.DB
}

func (p PostgresDB) Save(fullURL string) (string, error) {
	short := util.RandStringRunes()
	query := "insert into url_storager(short, long) values ($1, $2)"
	_, err := p.db.Exec(query, short, fullURL)
	if err != nil {
		p.log.Error("Failed to save short link into DB")
		return "", nil
	}
	p.log.Info("saved to db full url", fullURL)
	p.log.Info("short url is", short)
	return short, nil
}

func (p PostgresDB) Get(shortURL string) (string, error) {
	var long string
	query := "select long from url_storager where short=$1"
	row := p.db.QueryRow(query, shortURL)
	if err := row.Scan(&long); err != nil {
		p.log.Error("Failed to get link from db")
		return "", nil
	}
	return long, nil
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
		log.Error("Unable to open db, err is", err)
		return nil
	}
	_, err = db.Exec("create table IF NOT EXISTS url_storager(id serial, short text not null, long text not null)")
	if err != nil {
		log.Error("unable to create table, err is", err)
		return nil
	}
	return &PostgresDB{
		log: log,
		db:  db,
	}
}
