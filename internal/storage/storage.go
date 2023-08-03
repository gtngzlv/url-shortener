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
	SaveFull(userID string, fullURL string) (models.URLInfo, error)
	GetByShort(shortURL string) (models.URLInfo, error)
	GetBatchByUserID(userID string) ([]models.URLInfo, error)
	Batch(userID string, entities []models.URLInfo) ([]models.URLInfo, error)
	DeleteByUserIDAndShort(userID string, shortURL string) (bool, error)
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

func (s *storage) DeleteByUserIDAndShort(userID string, short string) (bool, error) {
	return s.defaultStorage.DeleteByUserIDAndShort(userID, short)
}

func (s *storage) Batch(userID string, entities []models.URLInfo) ([]models.URLInfo, error) {
	return s.defaultStorage.Batch(userID, entities)
}

func (s *storage) GetBatchByUserID(userID string) ([]models.URLInfo, error) {
	return s.defaultStorage.GetBatchByUserID(userID)
}

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

func (s *storage) GetByShort(short string) (models.URLInfo, error) {
	var urlInfo models.URLInfo
	full, err := s.defaultStorage.GetByShort(short)
	if err != nil {
		return urlInfo, err
	}
	return full, nil
}

func (s *storage) Ping() error {
	return nil
}
