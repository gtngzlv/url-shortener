package main

import (
	"github.com/gtngzlv/url-shortener/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
