package handlers

import (
	"github.com/gtngzlv/url-shortener/internal/pkg"
	"io"
	"log"
	"net/http"
)

func PostURL(w http.ResponseWriter, r *http.Request) {
	seps := []rune{';'}
	contentType := r.Header.Get("Content-Type")
	if pkg.SplitString(contentType, seps)[0] != "text/plain" {
		w.WriteHeader(400)
		log.Print("PostURL: incorrect format of content-type")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("PostURL: error: %s while reading body", err)
	}

	if len(string(body)) == 0 {
		log.Println("PostURL: Empty request body")
		return
	}

	shorted := pkg.GenerateShortURL(string(body))
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + shorted))
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("GetURL: err %s while parse form\n", err)
	}
	val := r.URL.Path

	longURL := pkg.GetFromStorage(val[1:])
	w.Header().Add("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
