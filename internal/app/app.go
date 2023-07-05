package app

import (
	"net/http"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"github.com/gtngzlv/url-shortener/internal/logger"
	"github.com/gtngzlv/url-shortener/internal/storage"
)

func Run() error {

	log := logger.NewLogger()
	cfg := config.LoadConfig()
	storage := storage.Init(cfg, log)
	app := handlers.NewApp(cfg, log, storage)
	return http.ListenAndServe(cfg.Host, app)
}
