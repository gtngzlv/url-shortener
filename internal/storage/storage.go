package storage

import (
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/errors"
	"github.com/gtngzlv/url-shortener/internal/models"
	"github.com/gtngzlv/url-shortener/internal/storage/database"
	"github.com/gtngzlv/url-shortener/internal/storage/filestorage"
)

type MyStorage interface {
	Save(fullURL string) (string, error)
	Get(shortURL string) (string, error)
	Batch(entities []models.BatchEntity) ([]models.BatchEntity, error)
	Ping() error
}

var Cache = make(map[string]string)

type storage struct {
	defaultStorage MyStorage
}

func (s *storage) Batch(entities []models.BatchEntity) ([]models.BatchEntity, error) {
	return s.defaultStorage.Batch(entities)
}

func Init(log zap.SugaredLogger, cfg *config.AppConfig) MyStorage {
	var s storage
	if cfg.DatabaseDSN != "" {
		s.defaultStorage = database.Init(log, cfg)
	} else if cfg.FileStoragePath != "" {
		s.defaultStorage = filestorage.Init(log, cfg.FileStoragePath)
	}
	return &s
}

func (s *storage) Save(full string) (string, error) {
	short, err := s.defaultStorage.Save(full)

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

func (s *storage) Get(short string) (string, error) {
	if getFromStorage(short) != "" {
		return getFromStorage(short), nil
	}
	full, err := s.defaultStorage.Get(short)
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
