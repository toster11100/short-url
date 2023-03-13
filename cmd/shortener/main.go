package main

import (
	"log"

	"github.com/toster11100/shortUrl.git/internal/server"
)

func main() {
	a := server.New()
	err := a.Start()
	if err != nil {
		log.Fatal(err)
	}
}
