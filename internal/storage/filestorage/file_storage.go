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

type fileStorage struct {
	path string
	log  zap.SugaredLogger
}

// Init inits file storage
func Init(log zap.SugaredLogger, fileStoragePath string) *fileStorage {
	return &fileStorage{
		path: fileStoragePath,
		log:  log,
	}
}

// GetBatchByUserID return batch of full urls by userID
func (f *fileStorage) GetBatchByUserID(userID string) ([]models.URLInfo, error) {
	file, err := os.OpenFile(f.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf("Get batch by userid: failed to close file, err: %s", err)
		}
	}(file)
	urls, err := readFromFileByUserID(file, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// Batch saves batch of urls
func (f *fileStorage) Batch(userID string, entities []models.URLInfo) ([]models.URLInfo, error) {
	var urls []models.URLInfo
	if _, err := os.Stat(filepath.Dir(f.path)); os.IsNotExist(err) {
		f.log.Infof("Creating folder")
		err = os.Mkdir(filepath.Dir(f.path), 0755)
		if err != nil {
			f.log.Infof("Error: %s", err)
		}
	}

	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		f.log.Infof("batch: error while OpenFile is %s\n", err)
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf("batch: failed to close file, err: %s", err)
		}
	}(file)
	f.log.Info("Created file by path", f.path)

	for _, v := range entities {
		shortURL := util.RandStringRunes()
		event := models.URLInfo{
			UserID:      userID,
			UUID:        uuid.NewString(),
			ShortURL:    shortURL,
			OriginalURL: v.OriginalURL,
		}
		data, err := json.Marshal(event)
		if err != nil {
			f.log.Infof("fileStorage Save: error while json marshal is %s", err)
			return nil, err
		}
		data = append(data, '\n')
		_, err = file.Write(data)
		if err != nil {
			f.log.Infof("fileStorage Save: error while write is %s", err)
			return nil, err
		}
		urls = append(urls, event)
	}
	return urls, nil
}

// SaveFull save full url and return short
func (f *fileStorage) SaveFull(userID string, fullURL string) (models.URLInfo, error) {
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
		f.log.Infof("fileStorage Save: error while OpenFile is %s\n", err)
		return urlInfo, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf("Save fileStorage: failed to close file, err: %s", err)
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
		f.log.Infof("fileStorage Save: error while json marshal is %s", err)
		return urlInfo, err
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		f.log.Infof("fileStorage Save: error while write is %s", err)
		return urlInfo, err
	}
	return event, nil
}

// GetByShort return full url by short
func (f *fileStorage) GetByShort(shortURL string) (models.URLInfo, error) {
	var urlInfo models.URLInfo
	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		f.log.Errorf("fileStorage Get: failed to get from file")
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf("Get fileStorage: failed to close file, err: %s", err)
		}
	}(file)

	if err != nil {
		f.log.Infof("fileStorage Get: error while OpenFile is %s", err)
		return urlInfo, err
	}
	url, err := readFromFileByShort(file, shortURL)
	if err != nil {
		f.log.Errorf("FS GetByShort: failed to get by short %s", err)
		return urlInfo, err
	}
	urlInfo.OriginalURL = url.OriginalURL
	return urlInfo, nil
}

// DeleteByUserIDAndShort marsk url as deleted in file storage
func (f *fileStorage) DeleteByUserIDAndShort(userID string, short string) error {
	file, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		f.log.Errorf("failed to get from file")
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			f.log.Errorf(" failed to close file, err: %s", err)
		}
	}(file)

	if err != nil {
		f.log.Infof("error while OpenFile is %s", err)
		return err
	}
	urlInfo, err := readFromFileByShort(file, short)
	if err != nil {
		f.log.Errorf("failed to get by short %s", err)
		return err
	}
	if urlInfo.UserID != userID {
		return err
	}
	urlInfo.IsDeleted = 1
	data, err := json.Marshal(urlInfo)
	if err != nil {
		f.log.Infof("error while json marshal is %s", err)
		return err
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		f.log.Infof("error while write is %s", err)
		return err
	}
	return nil
}

// Ping return nil
func (f *fileStorage) Ping() error {
	return nil
}

func readFromFileByShort(file *os.File, shortURL string) (models.URLInfo, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		var item models.URLInfo
		err := json.Unmarshal([]byte(text), &item)
		if err != nil {
			return item, err
		}
		if item.ShortURL == shortURL {
			return item, nil
		}
	}
	return models.URLInfo{}, nil
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
