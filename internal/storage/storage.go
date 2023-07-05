package storage

import (
	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/storage/filestorage"
	"go.uber.org/zap"
)

type MyStorage interface {
	Save(fullURL string) (string, error)
	Get(shortURL string) (string, error)
}

type storage struct {
	defaultStorage MyStorage
}

func Init(cfg *config.AppConfig, log zap.SugaredLogger) MyStorage {
	var s storage
	s.defaultStorage = filestorage.Init(log, cfg.FileStoragePath)
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
