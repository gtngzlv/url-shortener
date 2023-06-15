package filestorage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/pkg"
)

type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileStorage struct {
	Path string
}

var storage FileStorage

func Init(cfg *config.AppConfig) *FileStorage {
	storage = FileStorage{
		Path: cfg.FileStoragePath,
	}
	return &storage
}

func (f *FileStorage) Save(fullURL string) (string, error) {
	if shortURL := getShortURLFromStorage(fullURL); shortURL != "" {
		return shortURL, nil
	}

	if _, err := os.Stat(filepath.Dir(storage.Path)); os.IsNotExist(err) {
		log.Println("Creating folder")
		err = os.Mkdir(filepath.Dir(storage.Path), 0755)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}

	file, err := os.OpenFile(storage.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	log.Println("Created file by path", storage.Path)
	if err != nil {
		log.Fatalf("FileStorage Save: error while OpenFile is %s", err)
		return "", err
	}
	event := Event{
		UUID:        uuid.NewString(),
		ShortURL:    pkg.RandStringRunes(),
		OriginalURL: fullURL,
	}
	data, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("FileStorage Save: error while json marshal is %s", err)
		return "", err
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("FileStorage Save: error while write is %s", err)
		return "", nil
	}
	WriteToCache(fullURL, event.ShortURL)
	return event.ShortURL, nil
}

var Storage = make(map[string]string)

func (f *FileStorage) Get(shortURL string) (string, error) {
	file, err := os.OpenFile(storage.Path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("FileStorage Get: error while OpenFile is %s", err)
		return "", nil
	}
	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item Event
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			return "", nil
		}

		if item.ShortURL == shortURL {
			return item.OriginalURL, nil
		}
		s, e = readLine(r)
	}
	return "", nil
}

func getShortURLFromStorage(fullURL string) string {
	file, err := os.OpenFile(storage.Path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return ""
	}

	r := bufio.NewReader(file)
	s, e := readLine(r)
	var item Event
	for e == nil {
		err = json.Unmarshal([]byte(s), &item)
		if err != nil {
			return ""
		}

		if item.OriginalURL == fullURL {
			return item.ShortURL
		}
		s, e = readLine(r)
	}
	return ""
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix       = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
