package pkg

import (
	"math/rand"
	"strings"
)

var storage = make(map[string]string)

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
	return storage[key]
}

func GenerateShortURL(baseURL string) string {
	shortURL := randStringRunes()
	storage[shortURL] = baseURL
	return shortURL
}
