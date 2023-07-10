package database

import (
	"context"
	"time"

	"database/sql"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/errors"
	"github.com/gtngzlv/url-shortener/internal/models"
	"github.com/gtngzlv/url-shortener/internal/util"
)

type PostgresDB struct {
	log zap.SugaredLogger
	db  *sql.DB
	cfg *config.AppConfig
}

var tableName = "url_storage"

func (p PostgresDB) Batch(entities []models.BatchEntity) ([]models.BatchEntity, error) {
	var resultEntities []models.BatchEntity
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	tx, err := p.db.Begin()
	if err != nil {
		p.log.Error("Error while begin tx")
		return resultEntities, err
	}
	for _, v := range entities {
		short := util.RandStringRunes()
		_, err = tx.ExecContext(ctx, "INSERT INTO "+tableName+"(short, long) values ($1, $2)", short, v.OriginalURL)
		if err != nil {
			p.log.Error("Error while ExecContext", err)
			tx.Rollback()
			return resultEntities, nil
		}
		resultEntities = append(resultEntities, models.BatchEntity{
			CorrelationID: v.CorrelationID,
			ShortURL:      p.cfg.ResultURL + "/" + short,
		})
	}
	return resultEntities, tx.Commit()
}

func (p PostgresDB) Save(fullURL string) (string, error) {
	var (
		short string
		err   error
	)
	short = util.RandStringRunes()
	query := "INSERT INTO " + tableName + "(short, long) VALUES ($1, $2)"
	_, err = p.db.Exec(query, short, fullURL)
	if err != nil {
		p.log.Info("DB Save err ", err)
		if pgerrcode.IsIntegrityConstraintViolation(string(err.(*pq.Error).Code)) {
			short, err = p.GetShortURL(fullURL)
			if err != nil {
				p.log.Error("failed to get already saved short url")
				return "", nil
			}
			return short, errors.ErrAlreadyExist
		}
		p.log.Error("Failed to save short link into DB")
		return "", nil
	}
	return short, nil
}

func (p PostgresDB) Get(shortURL string) (string, error) {
	var long string
	query := "select long from " + tableName + " where short=$1"
	row := p.db.QueryRow(query, shortURL)
	if err := row.Scan(&long); err != nil {
		p.log.Error("Failed to get link from db")
		return "", nil
	}
	return long, nil
}

func (p PostgresDB) GetShortURL(fullURL string) (string, error) {
	var short string
	query := "select short from " + tableName + " where long=$1"
	row := p.db.QueryRow(query, fullURL)
	if err := row.Scan(&short); err != nil {
		return "", err
	}
	return short, nil
}

func (p PostgresDB) Ping() error {
	if err := p.db.Ping(); err != nil {
		return err
	}
	return nil
}

func Init(log zap.SugaredLogger, config *config.AppConfig) *PostgresDB {
	db, err := sql.Open("postgres", config.DatabaseDSN)
	if err != nil {
		log.Error("Unable to open db, err is", err)
		return nil
	}

	if err = goose.SetDialect("postgres"); err != nil {
		log.Error("unable to set goose dialect", err)
		return nil
	}
	if err = goose.Up(db, "migrations"); err != nil {
		log.Error("failed to load migrations ", err)
		return nil
	}
	return &PostgresDB{
		log: log,
		db:  db,
		cfg: config,
	}
}
