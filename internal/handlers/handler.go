package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/toster11100/shortUrl.git/internal/storage"
)

type Repositories interface {
	ReadUrl(int) string
	WriteUrl(string, int)
}

type Server struct {
	urlMap  Repositories
	id      int
	handler http.Handler
}

func New() *Server {
	mux := http.NewServeMux()
	MyServer := &Server{
		urlMap:  make(storage.Rep),
		handler: mux,
		id:      1,
	}
	mux.HandleFunc("/", MyServer.CheckMethod)
	return MyServer
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) CheckMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createShortURL(w, r)
		return
	case http.MethodGet:
		s.redirectToLongURL(w, r)
		return
	default:
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}
}

func (s *Server) createShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	strBody := strings.TrimRight(string(body), "\n")
	if _, err = url.ParseRequestURI(strBody); err != nil {
		http.Error(w, "this is not URL", http.StatusBadRequest)
		return
	}

	s.urlMap.WriteUrl(strBody, s.id)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://localhost:8080/"+strconv.Itoa(s.id))
	s.id++
}

func (s *Server) redirectToLongURL(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 || s.urlMap.ReadUrl(id) == "" {
		http.Error(w, "invalid short id", http.StatusBadRequest)
		return
	}

	longURL := s.urlMap.ReadUrl(id)

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
