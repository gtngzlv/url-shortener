package config

import (
	"flag"
	"os"
	"strings"
)

// env variables
var (
	ServerAddress   = "SERVER_ADDRESS"
	BaseURL         = "BASE_URL"
	FileStoragePath = "FILE_STORAGE_PATH"
	DatabaseDSN     = "DATABASE_DSN"
	EnableHTTPS     = "ENABLE_HTTPS"
)

// AppConfig contains environment variables which should be set
type AppConfig struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
	EnableHTTPS     string
}

// LoadConfig gets env vars from arguments or environment
func LoadConfig() *AppConfig {
	config := &AppConfig{}
	getArgs(config)
	getENVs(config)
	return config
}

func getArgs(cfg *AppConfig) {
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Default ServerAddress:port")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Default result URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db-7.json", "Default File Storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&cfg.EnableHTTPS, "s", "", "Server would be run on TLS")
	flag.Parse()
}

func getENVs(cfg *AppConfig) {
	envRunAddr := strings.TrimSpace(os.Getenv(ServerAddress))
	if envRunAddr != "" {
		cfg.ServerAddress = envRunAddr
	}

	envBaseURL := strings.TrimSpace(os.Getenv(BaseURL))
	if envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}

	fileStorageFile := strings.TrimSpace(os.Getenv(FileStoragePath))
	if fileStorageFile != "" {
		cfg.FileStoragePath = fileStorageFile
	}

	databaseDSN := strings.TrimSpace(os.Getenv(DatabaseDSN))
	if databaseDSN != "" {
		cfg.DatabaseDSN = databaseDSN
	}

	enableHTTPS := strings.TrimSpace(os.Getenv(EnableHTTPS))
	if enableHTTPS != "" {
		cfg.EnableHTTPS = enableHTTPS
	}
}
