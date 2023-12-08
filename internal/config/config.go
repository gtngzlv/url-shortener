package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// env variables
var (
	ServerAddress   = "SERVER_ADDRESS"
	BaseURL         = "BASE_URL"
	FileStoragePath = "FILE_STORAGE_PATH"
	DatabaseDSN     = "DATABASE_DSN"
	EnableHTTPS     = "ENABLE_HTTPS"
	TrustedSubnet   = "TRUSTED_SUBNET"
	Path            = "CONFIG"
)

// AppConfig contains environment variables which should be set
type AppConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	Path            string
	TrustedSubnet   string `json:"trusted_subnet"`
}

// LoadConfig gets env vars from arguments or environment
func LoadConfig() *AppConfig {
	config := &AppConfig{}
	getArgs(config)
	getENVs(config)
	if config.Path != "" {
		config = getConfigFile(config.Path)
	}
	return config
}

func getArgs(cfg *AppConfig) {
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Default ServerAddress:port")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Default result URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db-7.json", "Default File Storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database DSN")
	flag.BoolVar(&cfg.EnableHTTPS, "s", false, "Boolean flag to run server on https")
	flag.StringVar(&cfg.Path, "c", "", "Config path")
	flag.StringVar(&cfg.TrustedSubnet, "t", "", "Trusted subnet")
	flag.Parse()
}

func getENVs(cfg *AppConfig) {
	srvAddr := strings.TrimSpace(os.Getenv(ServerAddress))
	if srvAddr != "" {
		cfg.ServerAddress = srvAddr
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

	httpsVar, err := strconv.ParseBool(os.Getenv(EnableHTTPS))
	if err != nil {
		cfg.EnableHTTPS = false
	}
	cfg.EnableHTTPS = httpsVar

	confPath := strings.TrimSpace(os.Getenv(Path))
	if confPath != "" {
		cfg.Path = confPath
	}

	trustedSubnet := strings.TrimSpace(os.Getenv(TrustedSubnet))
	if trustedSubnet != "" {
		cfg.TrustedSubnet = trustedSubnet
	}
}

func getConfigFile(filename string) *AppConfig {
	var cfg *AppConfig
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Print("getConfigFile: failed to read flie", err)
		return nil
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		log.Print("getConfigFile: failed to unmarshal config", err)
		return nil
	}
	return cfg
}
