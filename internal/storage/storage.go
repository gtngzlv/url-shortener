package storage

import (
	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/database"
	"github.com/gtngzlv/url-shortener/internal/filestorage"
	"go.uber.org/zap"
)

type MyStorage interface {
	Save(fullURL string) (string, error)
	Get(shortURL string) (string, error)
	Ping() error
}

type storage struct {
	defaultStorage MyStorage
}

func Init(log zap.SugaredLogger, cfg *config.AppConfig) MyStorage {
	var s storage
	if cfg.DatabaseDSN != "" {
		s.defaultStorage = database.Init(log, cfg.DatabaseDSN)
	} else if cfg.FileStoragePath != "" {
		s.defaultStorage = filestorage.Init(log, cfg.FileStoragePath)
	}
	return &s
}

func (s *storage) Save(full string) (string, error) {
	short, err := s.defaultStorage.Save(full)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (s *storage) Get(short string) (string, error) {
	full, err := s.defaultStorage.Get(short)
	if err != nil {
		return "", err
	}
	return full, nil
}

func (s *storage) Ping() error {
	return nil
}
