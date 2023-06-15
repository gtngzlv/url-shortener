package storage

import (
	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/storage/filestorage"
)

type MyStorage interface {
	Save(fullURL string) (string, error)
	Get(shortURL string) (string, error)
}

var defaultStorage MyStorage

func Init(cfg *config.AppConfig) {
	defaultStorage = filestorage.Init(cfg)
}

func SaveURL(full string) (string, error) {
	short, err := defaultStorage.Save(full)
	if err != nil {
		return "", err
	}
	return short, nil
}

func GetFullURL(short string) (string, error) {
	full, err := defaultStorage.Get(short)
	if err != nil {
		return "", err
	}
	return full, nil
}
