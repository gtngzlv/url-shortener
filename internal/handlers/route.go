package handlers

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/gtngzlv/url-shortener/internal/storage"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/gzip"
	"github.com/gtngzlv/url-shortener/internal/logger"
)

type app struct {
	Router  *chi.Mux
	cfg     *config.AppConfig
	log     zap.SugaredLogger
	storage storage.MyStorage
}

// NewApp return object of new app
func NewApp(router *chi.Mux, cfg *config.AppConfig, log zap.SugaredLogger, s storage.MyStorage) *app {
	a := &app{
		router,
		cfg,
		log,
		s,
	}
	a.reg()
	return a
}

func (a *app) reg() {
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
	a.Router.Get("/ping", a.Ping)
	a.Router.Get("/api/user/urls", a.GetURLs)

	a.Router.Post("/", a.PostURL)
	a.Router.Post("/api/shorten/batch", a.Batch)

	a.Router.Delete("/api/user/urls", a.DeleteURLs)

	a.Router.HandleFunc("/debug/pprof/*", pprof.Index)
	a.Router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	a.Router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	a.Router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	a.Router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	a.Router.HandleFunc("/debug/vars", expVars)

	a.Router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	a.Router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	a.Router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	a.Router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	a.Router.Handle("/debug/pprof/block", pprof.Handler("block"))
	a.Router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
}

func expVars(w http.ResponseWriter, _ *http.Request) {
	first := true
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
