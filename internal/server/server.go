package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/toster11100/shortUrl.git/internal/handlers"
	"github.com/toster11100/shortUrl.git/internal/storage"
)

type Server struct {
	addr    string
	handler http.Handler
	srv     http.Server
}

func New() *Server {
	repo := storage.New()
	hand := handlers.New(repo)

	Server := &Server{
		addr:    ":8080",
		handler: hand,
		srv:     http.Server{},
	}

	return Server
}

func (s *Server) Start() error {
	log.Println("starting server")
	err := http.ListenAndServe(s.addr, s.handler)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	log.Println("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
