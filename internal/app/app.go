package app

import (
	"net/http"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"github.com/gtngzlv/url-shortener/internal/logger"
	"github.com/gtngzlv/url-shortener/internal/storage"
)

func Run() error {
	logger.NewLogger()
	cfg := config.LoadConfig()
	storage.Init(cfg)
	app := handlers.NewApp(cfg)
	return http.ListenAndServe(cfg.Host, app)
}
