package main

import (
	"log"

	"github.com/gtngzlv/url-shortener/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
