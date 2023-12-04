package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/errors"
	"github.com/gtngzlv/url-shortener/internal/models"
	"github.com/gtngzlv/url-shortener/internal/util"
)

type postgresDB struct {
	log       zap.SugaredLogger
	db        *sql.DB
	resultURL string
}

var tableName = "url_storage"

// GetStatistic - return num of saved urls and users
func (p postgresDB) GetStatistic() *models.Statistic {
	var st *models.Statistic
	query := "select count(distinct userID), count(*) from " + tableName
	res, err := p.db.Query(query)
	if err != nil {
		return nil
	}
	err = res.Scan(&st.Users, &st.URLs)
	if err != nil {
		return nil
	}
	return st
}

// Batch saves batch of urls and returns batch of short urls
func (p postgresDB) Batch(userID string, entities []models.URLInfo) ([]models.URLInfo, error) {
	var resultEntities []models.URLInfo
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*70)
	defer cancel()
	tx, err := p.db.Begin()
	if err != nil {
		tx.Rollback()
		p.log.Error("Error while begin tx")
		return resultEntities, err
	}
	for _, v := range entities {
		short := util.RandStringRunes()
		_, err = tx.ExecContext(ctx, "INSERT INTO "+tableName+"(short, long, userID) values ($1, $2, $3)", short, v.OriginalURL, userID)
		if err != nil {
			p.log.Error("Error while ExecContext", err)
			tx.Rollback()
			return resultEntities, nil
		}
		resultEntities = append(resultEntities, models.URLInfo{
			CorrelationID: v.CorrelationID,
			ShortURL:      p.resultURL + "/" + short,
		})
	}
	return resultEntities, tx.Commit()
}

// GetBatchByUserID returns batch of saved urls for user
func (p postgresDB) GetBatchByUserID(userID string) ([]models.URLInfo, error) {
	var (
		entity models.URLInfo
		result []models.URLInfo
	)
	query := "select short, long from " + tableName + " where userID=$1"
	rows, err := p.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		err = rows.Scan(&entity.ShortURL, &entity.OriginalURL)
		if err != nil {
			break
		}
		entity.ShortURL = p.resultURL + "/" + entity.ShortURL
		result = append(result, entity)
	}
	if len(result) == 0 {
		return nil, errors.ErrNoBatchByUserID
	}
	return result, nil
}

// SaveFull saves full error and return error if already exists and short url if not
func (p postgresDB) SaveFull(userID string, fullURL string) (models.URLInfo, error) {
	var (
		urlInfo models.URLInfo
		err     error
	)
	generatedShort := util.RandStringRunes()
	query := "INSERT INTO " + tableName + "(short, long, userID) VALUES ($1, $2, $3)"
	_, err = p.db.Exec(query, generatedShort, fullURL, userID)
	if err != nil {
		p.log.Info("DB Save err ", err)
		if pgerrcode.IsIntegrityConstraintViolation(string(err.(*pq.Error).Code)) {
			urlInfo, err = p.GetShortURL(fullURL)
			if err != nil {
				p.log.Error("failed to get already saved short url")
				return models.URLInfo{
					ShortURL: "",
				}, nil
			}
			return urlInfo, errors.ErrAlreadyExist
		}
		p.log.Error("Failed to save short link into DB")
		return models.URLInfo{
			ShortURL: "",
		}, nil
	}
	urlInfo.ShortURL = generatedShort
	return urlInfo, nil
}

// GetByShort returns full url by short url
func (p postgresDB) GetByShort(shortURL string) (models.URLInfo, error) {
	p.log.Infof("GetByShort: received url %s", shortURL)
	var urlEntity models.URLInfo
	query := "select userID, long, is_deleted from " + tableName + " where short=$1"
	row := p.db.QueryRow(query, shortURL)
	if err := row.Scan(&urlEntity.UserID, &urlEntity.OriginalURL, &urlEntity.IsDeleted); err != nil {
		p.log.Error("Failed to get link from db")
		return urlEntity, err
	}
	p.log.Infof("GetByShort found url info: %v", urlEntity)
	return urlEntity, nil
}

// DeleteByUserIDAndShort delete full url from db by userID and short url
func (p postgresDB) DeleteByUserIDAndShort(userID string, short string) error {
	query := "UPDATE " + tableName + " SET is_deleted=1::bit WHERE userID=$1 and short=$2"
	rows, err := p.db.Exec(query, userID, short)
	if err != nil {
		return err
	}
	if r, err := rows.RowsAffected(); err != nil || r != int64(1) {
		p.log.Infof("0 rows affected in delete")
		return err
	}
	p.log.Infof("Marked as deleted link %s", short)
	return nil
}

// GetShortURL returns short url from db
func (p postgresDB) GetShortURL(fullURL string) (models.URLInfo, error) {
	var urlInfo models.URLInfo
	query := "select short from " + tableName + " where long=$1"
	row := p.db.QueryRow(query, fullURL)
	if err := row.Scan(&urlInfo.ShortURL); err != nil {
		return urlInfo, err
	}
	return urlInfo, nil
}

// Ping ping db
func (p postgresDB) Ping() error {
	if err := p.db.Ping(); err != nil {
		return err
	}
	return nil
}

// Init inits sql db
func Init(log zap.SugaredLogger, db *sql.DB, resultURL string) *postgresDB {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Error("unable to set goose dialect", err)
		return nil
	}
	if err := goose.Up(db, "migrations"); err != nil {
		log.Error("failed to load migrations ", err)
		return nil
	}
	return &postgresDB{
		log:       log,
		db:        db,
		resultURL: resultURL,
	}
}
