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
}

// LoadConfig gets env vars from arguments or environment
func LoadConfig() *AppConfig {
	config := &AppConfig{}
	getArgs(config)
	getENVs(config)
	if config.Path != "" {
		getConfigFile(config, config.Path)
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
	flag.Parse()
}

func getENVs(cfg *AppConfig) {
	cfg.ServerAddress = returnEnvVar(ServerAddress)
	cfg.BaseURL = returnEnvVar(BaseURL)
	cfg.FileStoragePath = returnEnvVar(FileStoragePath)
	cfg.DatabaseDSN = returnEnvVar(DatabaseDSN)
	cfg.Path = returnEnvVar(Path)

	httpsVar, err := strconv.ParseBool(os.Getenv(EnableHTTPS))
	if err != nil {
		cfg.EnableHTTPS = false
	}
	cfg.EnableHTTPS = httpsVar
}

func getConfigFile(cfg *AppConfig, filename string) {
	file, err := os.ReadFile(cfg.Path)
	if err != nil {
		log.Print("getConfigFile: failed to read flie", err)
		return
	}
	err = json.Unmarshal(file, cfg)
	if err != nil {
		log.Print("getConfigFile: failed to unmarshal config", err)
		return
	}
}

func returnEnvVar(name string) string {
	variable := strings.TrimSpace(os.Getenv(name))
	if variable != "" {
		return variable
	}
	return ""
}
