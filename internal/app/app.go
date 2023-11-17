// Package app runs application
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"github.com/gtngzlv/url-shortener/internal/logger"
	"github.com/gtngzlv/url-shortener/internal/storage"
)

// Run starts application
func Run() {
	srv, err := runSrv()
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	sig := <-sigint
	log.Printf("Received signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown error: %v", err)
	}
}

func runSrv() (*http.Server, error) {
	router := chi.NewRouter()
	log := logger.NewLogger()
	cfg := config.LoadConfig()
	st := storage.Init(log, cfg)
	app := handlers.NewApp(router, cfg, log, st)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	if cfg.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.ServerAddress),
		}
		srv.TLSConfig = manager.TLSConfig()
		return srv, srv.ListenAndServeTLS("", "")
	}
	return srv, http.ListenAndServe(cfg.ServerAddress, app.Router)
}
