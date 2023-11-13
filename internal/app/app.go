// Package app runs application
package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"github.com/gtngzlv/url-shortener/internal/logger"
	"github.com/gtngzlv/url-shortener/internal/storage"
)

// Run starts application
func Run() error {
	router := chi.NewRouter()
	log := logger.NewLogger()
	cfg := config.LoadConfig()
	st := storage.Init(log, cfg)
	app := handlers.NewApp(router, cfg, log, st)
	if cfg.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.ServerAddress),
		}
		srv := &http.Server{
			Addr:      cfg.ServerAddress,
			Handler:   router,
			TLSConfig: manager.TLSConfig(),
		}
		return srv.ListenAndServeTLS("", "")
	}
	return http.ListenAndServe(cfg.ServerAddress, app.Router)

}
