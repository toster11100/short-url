package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/toster11100/shortUrl.git/internal/config"
	"github.com/toster11100/shortUrl.git/internal/handlers"
	"github.com/toster11100/shortUrl.git/internal/storage"
	storagepath "github.com/toster11100/shortUrl.git/internal/storage_path"
)

type Server struct {
	handler http.Handler
	srv     http.Server
	config  string
}

func New(config *config.Config) *Server {
	var repo handlers.Repositories
	if config.StoragePath != "" {
		repo = storagepath.New(config.StoragePath)
	} else {
		repo = storage.New()
	}

	hand := handlers.New(repo, config.BaseURL)

	Server := &Server{
		handler: hand,
		srv:     http.Server{},
		config:  config.Addr,
	}

	return Server
}

func (s *Server) Start() error {
	log.Printf("starting server to addres %s", s.config)
	err := http.ListenAndServe(s.config, handlers.GzipHandle(s.handler))
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
