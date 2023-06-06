package app

import (
	"github.com/gtngzlv/url-shortener/internal/handlers"
	"net/http"
)

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handlers.PostURL(w, r)
	} else if r.Method == http.MethodGet {
		handlers.GetURL(w, r)
	} else {
		w.WriteHeader(400)
	}
}

func Run() error {
	return http.ListenAndServe(":8080", http.HandlerFunc(webhook))
}
