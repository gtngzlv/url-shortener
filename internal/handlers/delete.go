package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gtngzlv/url-shortener/internal/core"
)

func (a *App) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	var inputArray []string
	userID, err := core.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urls, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Errorf("DeleteURLs: failed to read body, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(urls, &inputArray)
	if err != nil {
		a.log.Errorf("DeleteURLs: failed to unmarshal input request, %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputCh := addShortURLs(inputArray)
	go a.MarkAsDeleted(inputCh, userID)

	w.WriteHeader(http.StatusAccepted)
}

func (a *App) MarkAsDeleted(inputShort chan string, userID string) {
	for v := range inputShort {
		err := a.storage.DeleteByUserIDAndShort(userID, v)
		if err != nil {
			a.log.Infof("Failed to mark deleted by short %s", v)
		}
	}
}

func addShortURLs(input []string) chan string {
	inputCh := make(chan string, 10)

	go func() {
		defer close(inputCh)
		for _, url := range input {
			inputCh <- url
		}
	}()

	return inputCh
}
