package handlers

import (
	"encoding/json"
	"github.com/gtngzlv/url-shortener/internal/core"
	"io"
	"net/http"

	"github.com/gtngzlv/url-shortener/internal/errors"
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

	userID, err := core.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	url, err := a.storage.SaveFull(userID, request.URL)
	if err != nil {
		if err == errors.ErrAlreadyExist {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			response.Result = a.cfg.ResultURL + "/" + url.ShortURL
			res, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				a.log.Errorf("PostURL: error: %s while reading body", err)
				return
			}
			w.Write(res)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response.Result = a.cfg.ResultURL + "/" + url.ShortURL
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

	userID, err := core.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	url, err := a.storage.SaveFull(userID, string(body))
	if err != nil {
		if err == errors.ErrAlreadyExist {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(a.cfg.ResultURL + "/" + url.ShortURL))
			if err != nil {
				a.log.Errorf("PostURL: Failed to write in body")
			}
			return
		}
		a.log.Errorf("Failed to save URL in storage")
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(a.cfg.ResultURL + "/" + url.ShortURL))
	if err != nil {
		a.log.Errorf("PostURL: Failed to write in body")
	}
}
