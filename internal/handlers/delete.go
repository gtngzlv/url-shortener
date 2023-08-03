package handlers

import (
	"encoding/json"
	"github.com/gtngzlv/url-shortener/internal/core"
	"io"
	"net/http"
	"sync"
)

func (a *App) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	var inputArray []string
	userID, err := core.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	urls, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Errorf("DeleteURLs: failed to read body, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = json.Unmarshal(urls, &inputArray)
	if err != nil {
		a.log.Errorf("DeleteURLs: failed to unmarshal input request, %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	inputCh := addShortURLs(inputArray)
	go a.MarkAsDeleted(inputCh, userID, &wg)
	wg.Wait()

	w.WriteHeader(http.StatusAccepted)
}

func (a *App) MarkAsDeleted(inputShort chan string, userID string, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range inputShort {
		a.storage.DeleteByUserIDAndShort(userID, v)
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
