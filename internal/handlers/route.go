package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/gzip"
	"github.com/gtngzlv/url-shortener/internal/logger"
)

type App struct {
	*chi.Mux
	cfg *config.AppConfig
	log zap.SugaredLogger
}

func NewApp(cfg *config.AppConfig, log zap.SugaredLogger) *App {
	app := &App{
		chi.NewRouter(),
		cfg,
		log,
	}
	app.reg()
	return app
}

func (a *App) reg() {
	a.Use(middleware.Compress(5, "text/html",
		"application/x-gzip",
		"text/plain",
		"application/json"))
	a.Use(gzip.MiddlewareCompressGzip)
	a.Use(logger.WithLogging)

	a.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/shorten", a.PostAPIShorten)
	})

	a.Get("/{shortID}", a.GetURL)
	a.Post("/", a.PostURL)
}
