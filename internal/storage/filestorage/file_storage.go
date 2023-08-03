package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/models"
	"github.com/gtngzlv/url-shortener/internal/util"
)

type FileStorage struct {
	path string
	log  zap.SugaredLogger
}

func (f *FileStorage) DeleteByUserIDAndShort(userID string, short string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func Init(log zap.SugaredLogger, fileStoragePath string) *FileStorage {
	return &FileStorage{
		path: fileStoragePath,
		log:  log,
	}
}

func (f *FileStorage) GetBatchByUserID(userID string) ([]models.URLInfo, error) {
	file, err := os.OpenFile(f.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf("Get FileStorage: failed to close file, err: %s", err)
		}
	}(file)
	urls, err := readFromFileByUserID(file, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (f *FileStorage) Batch(userID string, entities []models.URLInfo) ([]models.URLInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FileStorage) SaveFull(userID string, fullURL string) (models.URLInfo, error) {
	var urlInfo models.URLInfo
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
		return urlInfo, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf("Save FileStorage: failed to close file, err: %s", err)
		}
	}(file)
	f.log.Info("Created file by path", f.path)

	event := models.URLInfo{
		UserID:      userID,
		UUID:        uuid.NewString(),
		ShortURL:    util.RandStringRunes(),
		OriginalURL: fullURL,
	}
	data, err := json.Marshal(event)
	if err != nil {
		f.log.Infof("FileStorage Save: error while json marshal is %s", err)
		return urlInfo, err
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		f.log.Infof("FileStorage Save: error while write is %s", err)
		return urlInfo, err
	}
	return event, nil
}

func (f *FileStorage) GetByShort(shortURL string) (models.URLInfo, error) {
	var urlInfo models.URLInfo
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
		return urlInfo, err
	}
	originalURL, err := readFromFileByShort(file, shortURL)
	if err != nil {
		f.log.Errorf("FS GetByShort: failed to get by short %s", err)
		return urlInfo, err
	}
	urlInfo.OriginalURL = originalURL
	return urlInfo, nil
}

func (f *FileStorage) Ping() error {
	return nil
}

func readFromFileByShort(file *os.File, shortURL string) (string, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		var item models.URLInfo
		err := json.Unmarshal([]byte(text), &item)
		if err != nil {
			return "", err
		}
		if item.ShortURL == shortURL {
			return item.OriginalURL, nil
		}
	}
	return "", nil
}

func readFromFileByUserID(file *os.File, userID string) ([]models.URLInfo, error) {
	var (
		item models.URLInfo
		urls []models.URLInfo
	)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		err := json.Unmarshal([]byte(text), &item)
		if err != nil {
			return nil, err
		}
		if item.UserID == userID {
			urls = append(urls, item)
		}
	}
	return urls, nil
}
