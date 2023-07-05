package storage

import (
	"github.com/gtngzlv/url-shortener/internal/filestorage"
)

type MyStorage interface {
	Save(fullURL string) (string, error)
	Get(shortURL string) (string, error)
}

type storage struct {
	defaultStorage MyStorage
}

func Init(fs *filestorage.FileStorage) MyStorage {
	var s storage
	s.defaultStorage = fs
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
