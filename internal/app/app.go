package app

import (
	"github.com/gtngzlv/url-shortener/internal/logger"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/handlers"
)

func Run() error {
	logger.NewLogger()
	router := chi.NewRouter()
	router.Get("/{value}", logger.WithLogging(http.HandlerFunc(handlers.GetURL)))
	router.Post("/api/shorten", logger.WithLogging(http.HandlerFunc(handlers.PostAPIShorten)))
	router.Post("/", logger.WithLogging(http.HandlerFunc(handlers.PostURL)))
	config.ParseAddresses()
	return http.ListenAndServe(config.GetSrvAddr(), router)
}
