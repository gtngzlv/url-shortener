package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gtngzlv/url-shortener/internal/storage"
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/gzip"
	"github.com/gtngzlv/url-shortener/internal/logger"
)

type App struct {
	Router  *chi.Mux
	cfg     *config.AppConfig
	log     zap.SugaredLogger
	storage storage.MyStorage
}

func NewApp(router *chi.Mux, cfg *config.AppConfig, log zap.SugaredLogger, s storage.MyStorage) *App {
	app := &App{
		router,
		cfg,
		log,
		s,
	}
	app.reg()
	return app
}

func (a *App) reg() {
	a.Router.Use(middleware.Compress(5, "text/html",
		"application/x-gzip",
		"text/plain",
		"application/json"))
	a.Router.Use(gzip.MiddlewareCompressGzip)
	a.Router.Use(logger.WithLogging)

	a.Router.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/shorten", a.PostAPIShorten)
	})

	a.Router.Get("/{shortID}", a.GetURL)
	a.Router.Post("/", a.PostURL)
}
