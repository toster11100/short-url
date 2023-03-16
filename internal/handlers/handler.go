package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Repositories interface {
	ReadURL(int) (string, error)
	WriteURL(string) int
}

type Server struct {
	urlMap  Repositories
	handler http.Handler
}

func New(storage Repositories) *Server {
	router := mux.NewRouter()
	myServer := &Server{
		urlMap:  storage,
		handler: router,
	}
	router.HandleFunc("/", myServer.createShortURL).Methods(http.MethodPost)
	router.HandleFunc("/{id:[0-9]+}", myServer.redirectToLongURL).Methods(http.MethodGet)

	return myServer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) createShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading equest body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	strBody := strings.TrimRight(string(body), "\n")
	if _, err = url.ParseRequestURI(strBody); err != nil {
		err := fmt.Errorf("this \"%s\" is not URL", strBody)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := s.urlMap.WriteURL(strBody)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://localhost:8080/", id)
}

func (s *Server) redirectToLongURL(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil {
		err := fmt.Errorf("this ID: %s is not valid", path)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	longURL, err := s.urlMap.ReadURL(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
