package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gtngzlv/url-shortener/internal/core"
	"github.com/gtngzlv/url-shortener/internal/errors"
	"net/http"
)

func (a *App) GetURL(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "shortID")
	longURL, err := a.storage.GetByShort(val)
	a.log.Infof("Found %s url by short %s", longURL, val)
	if err != nil {
		a.log.Errorf("Error while GetURL: %s", err)
	}
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *App) GetURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := core.GetUserToken(w, r)
	a.log.Infof("received userID from cookie %s", userID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	batch, err := a.storage.GetBatchByUserID(userID)
	if err == errors.ErrNoBatchByUserID {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := json.Marshal(batch)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
