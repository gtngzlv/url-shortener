package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"github.com/gtngzlv/url-shortener/internal/logger"
	"github.com/gtngzlv/url-shortener/internal/storage"
)

func Run() error {
	router := chi.NewRouter()
	log := logger.NewLogger()
	cfg := config.LoadConfig()
	st := storage.Init(log, cfg)
	app := handlers.NewApp(router, cfg, log, st)
	return http.ListenAndServe(cfg.Host, app.Router)
}
