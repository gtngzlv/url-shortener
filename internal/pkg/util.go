package pkg

import (
	"github.com/gtngzlv/url-shortener/internal/storage"
	"math/rand"
	"strings"
)

func randStringRunes() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SplitString(s string, separators []rune) []string {
	f := func(r rune) bool {
		for _, s := range separators {
			if r == s {
				return true
			}
		}
		return false
	}
	return strings.FieldsFunc(s, f)
}

func GetFromStorage(key string) string {
	return storage.Storage[key]
}

func SetShortURL(baseURL string) string {
	shortURL := randStringRunes()
	storage.Storage[shortURL] = baseURL
	return shortURL
}
