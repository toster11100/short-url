package handlers

import (
	"encoding/json"
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
	WriteURL(string) string
}

type Server struct {
	urlMap  Repositories
	handler http.Handler
}

type ShortJSON struct {
	URL    string `json:"url,omitempty"`
	Result string `json:"result,omitempty"`
}

func New(storage Repositories) *Server {
	router := mux.NewRouter()
	myServer := &Server{
		urlMap:  storage,
		handler: router,
	}
	router.HandleFunc("/", myServer.createShortURL).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", myServer.shortenJSON).Methods(http.MethodPost)
	router.HandleFunc("/{id:[0-9]+}", myServer.redirectToLongURL).Methods(http.MethodGet)

	return myServer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) createShortURL(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading equest body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	strBody := strings.TrimRight(string(requestBody), "\n")
	if _, err = url.ParseRequestURI(strBody); err != nil {
		err := fmt.Errorf("this \"%s\" is not URL", strBody)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortenedURL := s.urlMap.WriteURL(strBody)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, shortenedURL)
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

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) shortenJSON(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading equest body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	requestBodyJSON := ShortJSON{}

	if err := json.Unmarshal(requestBody, &requestBodyJSON); err != nil {
		log.Println("error unmarshaling JSON", err)
		http.Error(w, "invalid data JSON", http.StatusBadRequest)
		return
	}

	if requestBodyJSON.URL == "" {
		log.Println("URL field is empty")
		http.Error(w, "URL field is required", http.StatusBadRequest)
		return
	}

	shortenedURL := s.urlMap.WriteURL(requestBodyJSON.URL)

	requestBodyJSON = ShortJSON{
		Result: shortenedURL,
	}
	responseBody, err := json.Marshal(requestBodyJSON)
	if err != nil {
		log.Println("error marshaling JSON", err)
		http.Error(w, "invalid data JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)
}
