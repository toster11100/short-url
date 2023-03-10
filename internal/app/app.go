package app

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type server struct {
	urlMap  map[int]string
	addr    string
	handler http.Handler
	count   int
}

func Mew() *server {
	mux := http.NewServeMux()

	myServer := &server{
		urlMap:  make(map[int]string),
		addr:    ":8080",
		handler: mux,
		count:   1,
	}
	mux.HandleFunc("/", myServer.checkMethod)
	return myServer
}

func (s *server) checkMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createShortURL(w, r)
		return
	case http.MethodGet:
		s.redirectToLongURL(w, r)
		return
	default:
		http.Error(w, "invalid method ", http.StatusBadRequest)
		return
	}
}

func (s *server) createShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	strBody := strings.TrimRight(string(body), "\n")
	if len(strBody) == 0 {
		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	s.urlMap[s.count] = strBody

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "http://localhost:8080/"+strconv.Itoa(s.count))
	s.count++
}

func (s *server) redirectToLongURL(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 || s.urlMap[id] == "" {
		http.Error(w, "invalid short id", http.StatusBadRequest)
		return
	}

	longURL := s.urlMap[id]
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s server) Start() error {
	err := http.ListenAndServe(s.addr, s.handler)
	if err != nil {
		return err
	}
	return nil
}
