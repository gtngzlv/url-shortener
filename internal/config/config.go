package config

import (
	"flag"
	"os"
	"strings"
)

var (
	ServerAddress   = "SERVER_ADDRESS"
	BaseURL         = "BASE_URL"
	FileStoragePath = "FILE_STORAGE_PATH"
	DatabaseDSN     = "DATABASE_DSN"
)

type AppConfig struct {
	Host            string
	ResultURL       string
	FileStoragePath string
	DatabaseDSN     string
}

func LoadConfig() *AppConfig {
	config := &AppConfig{}
	getArgs(config)
	getENVs(config)
	return config
}

func getArgs(cfg *AppConfig) {
	flag.StringVar(&cfg.Host, "a", "localhost:8080", "Default Host:port")
	flag.StringVar(&cfg.ResultURL, "b", "http://localhost:8080", "Default result URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db-7.json", "Default File Storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database DSN")
	flag.Parse()
}

func getENVs(cfg *AppConfig) {
	envRunAddr := strings.TrimSpace(os.Getenv(ServerAddress))
	if envRunAddr != "" {
		cfg.Host = envRunAddr
	}

	envBaseURL := strings.TrimSpace(os.Getenv(BaseURL))
	if envBaseURL != "" {
		cfg.ResultURL = envBaseURL
	}

	fileStorageFile := strings.TrimSpace(os.Getenv(FileStoragePath))
	if fileStorageFile != "" {
		cfg.FileStoragePath = fileStorageFile
	}

	databaseDSN := strings.TrimSpace(os.Getenv(DatabaseDSN))
	if databaseDSN != "" {
		cfg.DatabaseDSN = databaseDSN
	}
}
