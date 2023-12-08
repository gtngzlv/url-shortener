package storage

import (
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/core"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/errors"
	"github.com/gtngzlv/url-shortener/internal/models"
	"github.com/gtngzlv/url-shortener/internal/storage/database"
	"github.com/gtngzlv/url-shortener/internal/storage/filestorage"
)

// MyStorage interface with methods to storage
type MyStorage interface {
	SaveFull(userID string, fullURL string) (models.URLInfo, error)
	GetByShort(shortURL string) (models.URLInfo, error)
	GetBatchByUserID(userID string) ([]models.URLInfo, error)
	Batch(userID string, entities []models.URLInfo) ([]models.URLInfo, error)
	DeleteByUserIDAndShort(userID string, shortURL string) error
	Ping() error

	GetStatistic() *models.Statistic
}

type storage struct {
	defaultStorage MyStorage
}

// Init inits storage
func Init(log zap.SugaredLogger, cfg *config.AppConfig) MyStorage {
	var s storage
	if cfg.DatabaseDSN != "" {
		db, resultURL := core.InitDB(cfg.DatabaseDSN, cfg.BaseURL)
		if db == nil {
			log.Error("Failed to init DB")
		}
		s.defaultStorage = database.Init(log, db, resultURL)
	} else if cfg.FileStoragePath != "" {
		s.defaultStorage = filestorage.Init(log, cfg.FileStoragePath)
	}
	return &s
}

// DeleteByUserIDAndShort marks url as deleted in storage
func (s *storage) DeleteByUserIDAndShort(userID string, short string) error {
	return s.defaultStorage.DeleteByUserIDAndShort(userID, short)
}

// Batch saves batch of urls in storage
func (s *storage) Batch(userID string, entities []models.URLInfo) ([]models.URLInfo, error) {
	return s.defaultStorage.Batch(userID, entities)
}

// GetBatchByUserID return batch of urls for user by id
func (s *storage) GetBatchByUserID(userID string) ([]models.URLInfo, error) {
	return s.defaultStorage.GetBatchByUserID(userID)
}

// SaveFull saves full url in storage, returns short and err if exist
func (s *storage) SaveFull(userID string, full string) (models.URLInfo, error) {
	var urlInfo models.URLInfo
	urlInfo, err := s.defaultStorage.SaveFull(userID, full)

	switch {
	case err == errors.ErrAlreadyExist:
		{
			return urlInfo, err
		}
	case err != nil && err != errors.ErrAlreadyExist:
		{
			return urlInfo, err
		}
	default:
		{
			return urlInfo, nil
		}
	}
}

// GetByShort return full url by short
func (s *storage) GetByShort(short string) (models.URLInfo, error) {
	var urlInfo models.URLInfo
	full, err := s.defaultStorage.GetByShort(short)
	if err != nil {
		return urlInfo, err
	}
	return full, nil
}

// Ping returns nil to ping
func (s *storage) Ping() error {
	return nil
}

// GetStatistic - returns num of saved urls and users
func (s *storage) GetStatistic() *models.Statistic {
	return s.defaultStorage.GetStatistic()
}
