package handlers

import "net/http"

func (a *App) Ping(w http.ResponseWriter, r *http.Request) {
	if err := a.storage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
