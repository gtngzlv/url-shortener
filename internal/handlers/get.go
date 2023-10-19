package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gtngzlv/url-shortener/internal/core"
	"github.com/gtngzlv/url-shortener/internal/errors"
)

// GetURL returns full url by provided shortID
func (a *app) GetURL(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "shortID")
	url, err := a.storage.GetByShort(val)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	a.log.Infof("Found %s url by short %s", url.OriginalURL, val)
	if url.IsDeleted == 1 {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// GetURLs returns array of all saved by user urls
func (a *app) GetURLs(w http.ResponseWriter, r *http.Request) {
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
