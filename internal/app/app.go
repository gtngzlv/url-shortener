package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/gzip"
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"github.com/gtngzlv/url-shortener/internal/logger"
)

func Run() error {
	logger.NewLogger()
	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "text/html",
		"application/x-gzip",
		"text/plain",
		"application/json"))
	r.Use(gzip.Middleware)
	r.Use(logger.WithLogging)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/shorten", handlers.PostAPIShorten)
	})

	r.Get("/{shortID}", handlers.GetURL)
	r.Post("/", handlers.PostURL)
	config.ParseAddresses()
	return http.ListenAndServe(config.GetSrvAddr(), r)
}
