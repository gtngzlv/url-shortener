package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/util"
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
	if _, err := os.Stat(filepath.Dir(f.path)); os.IsNotExist(err) {
		f.log.Infof("Creating folder")
		err = os.Mkdir(filepath.Dir(f.path), 0755)
		if err != nil {
			f.log.Infof("Error: %s", err)
		}
	}

	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		f.log.Infof("FileStorage Save: error while OpenFile is %s\n", err)
		return "", err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf("Save FileStorage: failed to close file, err: %s", err)
		}
	}(file)
	f.log.Info("Created file by path", f.path)

	event := Event{
		UUID:        uuid.NewString(),
		ShortURL:    util.RandStringRunes(),
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
	Cache[event.ShortURL] = fullURL
	return event.ShortURL, nil
}

func (f *FileStorage) Get(shortURL string) (string, error) {
	if Cache[shortURL] != "" {
		return Cache[shortURL], nil
	}
	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		f.log.Errorf("FileStorage Get: failed to get from file")
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf("Get FileStorage: failed to close file, err: %s", err)
		}
	}(file)

	if err != nil {
		f.log.Infof("FileStorage Get: error while OpenFile is %s", err)
		return "", nil
	}
	return readFromFile(file, shortURL)
}

var Cache = make(map[string]string)

func readFromFile(file *os.File, shortURL string) (string, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		var item Event
		err := json.Unmarshal([]byte(text), &item)
		if err != nil {
			return "", err
		}
		if item.ShortURL == shortURL {
			return item.OriginalURL, nil
		}
		Cache[shortURL] = item.OriginalURL
	}
	return "", nil
}
