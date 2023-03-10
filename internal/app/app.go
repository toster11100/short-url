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
	Addr    string
	Handler http.Handler
	count   int
}

func Mew() *server {
	mux := http.NewServeMux()

	myServer := &server{
		urlMap:  make(map[int]string),
		Addr:    ":8080",
		Handler: mux,
		count:   1,
	}
	mux.HandleFunc("/", myServer.CheckMethod)
	return myServer
}

func (s *server) CheckMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.CreateShortURL(w, r)
		return
	case http.MethodGet:
		s.RedirectToLongURL(w, r)
		return
	default:
		http.Error(w, "invalid method ", http.StatusBadRequest)
		return
	}
}

func (s *server) CreateShortURL(w http.ResponseWriter, r *http.Request) {
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

func (s *server) RedirectToLongURL(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 || s.urlMap[id] == "" {
		http.Error(w, "invalid short id", http.StatusBadRequest)
		return
	}

	longUrl := s.urlMap[id]
	w.Header().Set("Location", longUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s server) Start() error {
	err := http.ListenAndServe(s.Addr, s.Handler)
	if err != nil {
		return err
	}
	return nil
}
