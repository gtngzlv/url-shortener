package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/gtngzlv/url-shortener/internal/filestorage"
	"net/http"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"github.com/gtngzlv/url-shortener/internal/logger"
	"github.com/gtngzlv/url-shortener/internal/storage"
)

func Run() error {
	router := chi.NewRouter()
	log := logger.NewLogger()
	cfg := config.LoadConfig()
	fileStorage := filestorage.Init(log, cfg.FileStoragePath)
	storage := storage.Init(fileStorage)
	app := handlers.NewApp(router, cfg, log, storage)
	return http.ListenAndServe(cfg.Host, app.Router)
}
