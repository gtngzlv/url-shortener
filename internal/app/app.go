package app

import (
	"github.com/gtngzlv/url-shortener/internal/handler"
	"net/http"
)

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handler.PostURL(w, r)
	} else if r.Method == http.MethodGet {
		handler.GetURL(w, r)
	} else {
		w.WriteHeader(400)
	}
}

func Run() error {
	return http.ListenAndServe(":8080", http.HandlerFunc(webhook))
}
