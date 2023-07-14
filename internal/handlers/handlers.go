package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gtngzlv/url-shortener/internal/models"
)

func (a *App) PostAPIShorten(w http.ResponseWriter, r *http.Request) {
	var (
		request  models.APIShortenRequest
		response models.APIShortenResponse
		err      error
	)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.log.Errorf("PostURL: error: %s while reading body", err)
		return
	}
	err = json.Unmarshal(bytes, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.log.Errorf("PostURL: error: %s while reading body", err)
		return
	}

	shorted, err := a.storage.Save(request.URL)
	a.log.Infof("Saved URL: %s", shorted)
	if err != nil {
		a.log.Errorf("Error while saving json short url: %s", err)
	}
	response.Result = a.cfg.ResultURL + "/" + shorted
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.log.Errorf("PostURL: error: %s while reading body", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (a *App) PostURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.log.Errorf("PostURL: error: %s while reading body", err)
		return
	}
	a.log.Infof("Received URL: %s", string(body))

	if len(string(body)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		a.log.Errorf("PostURL: Empty request body")
		return
	}

	shorted, err := a.storage.Save(string(body))
	if err != nil {
		a.log.Errorf("Failed to save URL in storage")
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(a.cfg.ResultURL + "/" + shorted))
	if err != nil {
		a.log.Errorf("PostURL: Failed to write in body")
	}
}

func (a *App) GetURL(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "shortID")
	longURL, err := a.storage.Get(val)
	a.log.Infof("Found %s url by short %s", longURL, val)
	if err != nil {
		a.log.Errorf("Error while GetURL: %s", err)
	}
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *App) Ping(w http.ResponseWriter, r *http.Request) {
	if err := a.storage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) Batch(w http.ResponseWriter, r *http.Request) {
	var batches []models.BatchEntity
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
	result, err := a.storage.Batch(batches)
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
