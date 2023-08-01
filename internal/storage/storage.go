package storage

import (
	"github.com/gtngzlv/url-shortener/internal/core"
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/errors"
	"github.com/gtngzlv/url-shortener/internal/models"
	"github.com/gtngzlv/url-shortener/internal/storage/database"
	"github.com/gtngzlv/url-shortener/internal/storage/filestorage"
)

type MyStorage interface {
	SaveFull(userID string, fullURL string) (string, error)
	GetByShort(shortURL string) (string, error)
	GetBatchByUserID(userID string) ([]models.BatchEntity, error)
	Batch(userID string, entities []models.BatchEntity) ([]models.BatchEntity, error)
	Ping() error
}

var Cache = make(map[string]string)

type storage struct {
	defaultStorage MyStorage
}

func Init(log zap.SugaredLogger, cfg *config.AppConfig) MyStorage {
	var s storage
	if cfg.DatabaseDSN != "" {
		db, resultURL := core.InitDB(cfg.DatabaseDSN, cfg.ResultURL)
		if db == nil {
			log.Error("Failed to init DB")
		}
		s.defaultStorage = database.Init(log, db, resultURL)
	} else if cfg.FileStoragePath != "" {
		s.defaultStorage = filestorage.Init(log, cfg.FileStoragePath)
	}
	return &s
}

func (s *storage) Batch(userID string, entities []models.BatchEntity) ([]models.BatchEntity, error) {
	return s.defaultStorage.Batch(userID, entities)
}

func (s *storage) GetBatchByUserID(userID string) ([]models.BatchEntity, error) {
	return s.defaultStorage.GetBatchByUserID(userID)
}

func (s *storage) SaveFull(userID string, full string) (string, error) {
	short, err := s.defaultStorage.SaveFull(userID, full)

	switch {
	case err == errors.ErrAlreadyExist:
		{
			return short, err
		}
	case err != nil && err != errors.ErrAlreadyExist:
		{
			return "", err
		}
	default:
		{
			saveToStorage(short, full)
			return short, nil
		}
	}
}

func (s *storage) GetByShort(short string) (string, error) {
	if getFromStorage(short) != "" {
		return getFromStorage(short), nil
	}
	full, err := s.defaultStorage.GetByShort(short)
	if err != nil {
		return "", err
	}
	return full, nil
}

func (s *storage) Ping() error {
	return nil
}

func saveToStorage(short, full string) {
	Cache[short] = full
}

func getFromStorage(short string) string {
	return Cache[short]
}
