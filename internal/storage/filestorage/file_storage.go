package filestorage

import (
	"bufio"
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/gtngzlv/url-shortener/internal/pkg"
)

type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileStorage struct {
	path string
	log  zap.SugaredLogger
}

func Init(log zap.SugaredLogger, fileStoragePath string) *FileStorage {
	return &FileStorage{
		path: fileStoragePath,
		log:  log,
	}
}

func (f *FileStorage) Save(fullURL string) (string, error) {
	if shortURL := f.getShortURLFromStorage(fullURL); shortURL != "" {
		return shortURL, nil
	}

	if _, err := os.Stat(filepath.Dir(f.path)); os.IsNotExist(err) {
		f.log.Infof("Creating folder")
		err = os.Mkdir(filepath.Dir(f.path), 0755)
		if err != nil {
			f.log.Infof("Error: %s", err)
		}
	}

	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			f.log.Errorf("Save FileStorage: failed to close file, err: %s", err)
		}
	}(file)

	f.log.Info("Created file by path", f.path)
	if err != nil {
		f.log.Infof("FileStorage Save: error while OpenFile is %s\n", err)
		return "", err
	}
	event := Event{
		UUID:        uuid.NewString(),
		ShortURL:    pkg.RandStringRunes(),
		OriginalURL: fullURL,
	}
	data, err := json.Marshal(event)
	if err != nil {
		f.log.Infof("FileStorage Save: error while json marshal is %s", err)
		return "", err
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		f.log.Infof("FileStorage Save: error while write is %s", err)
		return "", nil
	}
	WriteToCache(fullURL, event.ShortURL)
	return event.ShortURL, nil
}

func (f *FileStorage) Get(shortURL string) (string, error) {
	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0666)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			f.log.Errorf("Get FileStorage: failed to close file, err: %s", err)
		}
	}(file)

	if err != nil {
		f.log.Infof("FileStorage Get: error while OpenFile is %s", err)
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

func (f *FileStorage) getShortURLFromStorage(fullURL string) string {
	file, err := os.OpenFile(f.path, os.O_RDONLY|os.O_CREATE, 0666)
	defer file.Close()

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
