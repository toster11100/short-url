package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/toster11100/shortUrl.git/internal/config"
	"github.com/toster11100/shortUrl.git/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("error parsing addres %s", err)
		return
	}
	a := server.New(cfg)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.Start(); err != nil {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	log.Printf("received signal %s, shutting down", <-sigs)

	if err := a.Stop(); err != nil {
		log.Fatalf("server failed to shut down: %v", err)
		return
	}
}
