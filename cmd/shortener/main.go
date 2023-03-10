package main

import (
	"log"

	"github.com/toster11100/shortUrl.git/internal/app"
)

func main() {
	a := app.Mew()
	err := a.Start()
	if err != nil {
		log.Fatal(err)
	}
}
