package storage

import "github.com/gtngzlv/url-shortener/internal/pkg"

var Storage = make(map[string]string)

func GetValueFromStorage(key string) string {
	return Storage[key]
}

func ExistValueInStorage(value string) bool {
	for _, v := range Storage {
		if v == value {
			return true
		}
	}
	return false
}

func SetShortURL(baseURL string) string {
	shortURL := pkg.RandStringRunes()
	Storage[shortURL] = baseURL
	return shortURL
}
