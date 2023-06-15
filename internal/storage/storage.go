package storage

import "github.com/gtngzlv/url-shortener/internal/pkg"

var Storage = make(map[string]string)

func GetFromStorage(key string) string {
	return Storage[key]
}

func SetShortURL(baseURL string) string {
	shortURL := pkg.RandStringRunes()
	Storage[shortURL] = baseURL
	return shortURL
}
