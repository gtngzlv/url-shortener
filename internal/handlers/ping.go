package handlers

import "net/http"

// Ping makes test connection to storage
func (a *app) Ping(w http.ResponseWriter, r *http.Request) {
	if err := a.storage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
