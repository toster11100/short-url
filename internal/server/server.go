package server

import (
	"net/http"

	"github.com/toster11100/shortUrl.git/internal/handlers"
)

type server struct {
	addr    string
	handler http.Handler
}

func New() *server {
	hand := handlers.New()
	Server := &server{
		addr:    ":8080",
		handler: hand,
	}
	return Server
}

func (s server) Start() error {
	err := http.ListenAndServe(s.addr, s.handler)
	if err != nil {
		return err
	}
	return nil
}
