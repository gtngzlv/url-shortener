package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/pkg"
)

func PostURL(w http.ResponseWriter, r *http.Request) {
	seps := []rune{';'}
	contentType := r.Header.Get("Content-Type")
	if pkg.SplitString(contentType, seps)[0] != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("PostURL: incorrect format of content-type")
		return
	}

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

	shorted := pkg.SetShortURL(string(body))
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	finAddr := config.GetFinAddr()
	w.Write([]byte(finAddr + "/" + shorted))
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("GetURL: err %s while parse form\n", err)
		return
	}
	val := r.URL.Path

	longURL := pkg.GetFromStorage(val[1:])
	w.Header().Add("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
