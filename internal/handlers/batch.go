package handlers

import (
	"encoding/json"
	"github.com/gtngzlv/url-shortener/internal/core"
	"github.com/gtngzlv/url-shortener/internal/models"
	"io"
	"net/http"
)

func (a *App) Batch(w http.ResponseWriter, r *http.Request) {
	userID, err := core.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var batches []models.URLInfo
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("Batch: failed to read from body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &batches)
	a.log.Info("Batch request body", batches)
	if err != nil {
		a.log.Error("Batch: failed to unmarshal request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := a.storage.Batch(userID, batches)
	if err != nil {
		a.log.Error("Batch: failed to save to database")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(result)
	a.log.Info("Batch response", string(response))
	if err != nil {
		a.log.Error("Batch: failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
