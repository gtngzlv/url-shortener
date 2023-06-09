package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"net/http"
)

func Run() error {
	router := chi.NewRouter()
	router.Get("/{value}", handlers.GetURL)
	router.Post("/", handlers.PostURL)
	config.ParseAddresses()
	return http.ListenAndServe(config.GetSrvAddr(), router)
}
