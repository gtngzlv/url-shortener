package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gtngzlv/url-shortener/internal/storage"
	"io"
	"log"
	"net/http"

	"github.com/gtngzlv/url-shortener/internal/config"
)

func PostAPIShorten(w http.ResponseWriter, r *http.Request) {
	var (
		request  APIShortenRequest
		response APIShortenResponse
		err      error
	)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("PostURL: error: %s while reading body", err)
		return
	}
	err = json.Unmarshal(bytes, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("PostURL: error: %s while reading body", err)
		return
	}

	shorted := storage.SetShortURL(request.URL)
	finAddr := config.GetFinAddr()
	response.Result = finAddr + "/" + shorted
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("PostURL: error: %s while reading body", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func PostURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("PostURL: error: %s while reading body", err)
		return
	}

	if len(string(body)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("PostURL: Empty request body")
		return
	}

	shorted := storage.SetShortURL(string(body))
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	finAddr := config.GetFinAddr()

	w.Write([]byte(finAddr + "/" + shorted))
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "shortID")

	longURL := storage.GetValueFromStorage(val)
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
