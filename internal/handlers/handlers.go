package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *App) PostAPIShorten(w http.ResponseWriter, r *http.Request) {
	var (
		request  APIShortenRequest
		response APIShortenResponse
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
	if err != nil {
		a.log.Errorf("Error while GetURL: %s", err)
	}
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
