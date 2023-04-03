package handlers

import (
	"compress/gzip"
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
	WriteURL(string) (int, error)
}

type Server struct {
	urlMap  Repositories
	handler http.Handler
	cfg     string
}

type ShortJSON struct {
	URL    string `json:"url,omitempty"`
	Result string `json:"result,omitempty"`
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func New(storage Repositories, config string) *Server {
	router := mux.NewRouter()
	myServer := &Server{
		urlMap:  storage,
		handler: router,
		cfg:     config,
	}
	router.HandleFunc("/", myServer.createShortURL).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", myServer.shortenJSON).Methods(http.MethodPost)
	router.HandleFunc("/{id:[0-9]+}", myServer.redirectToLongURL).Methods(http.MethodGet)

	return myServer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandle(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Println("failed gzip reader:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer gzr.Close()
			r.Body = http.MaxBytesReader(w, gzr, r.ContentLength)
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			log.Println("failed gzip writer:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		gw := &gzipWriter{
			ResponseWriter: w,
			Writer:         gz,
		}
		handler.ServeHTTP(gw, r)
	})
}

func (s *Server) createShortURL(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading equest body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	strBody := string(requestBody)
	if _, err = url.ParseRequestURI(strBody); err != nil {
		err := fmt.Errorf("this \"%s\" is not URL", strBody)
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.urlMap.WriteURL(strBody)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortenedURL := fmt.Sprintf("%v/%v", s.cfg, id)

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
		log.Println("error unmarshaling JSON:", err)
		http.Error(w, "invalid data JSON", http.StatusBadRequest)
		return
	}

	if requestBodyJSON.URL == "" {
		log.Println("URL field is empty")
		http.Error(w, "URL field is required", http.StatusBadRequest)
		return
	}

	shortenedURL, err := s.urlMap.WriteURL(requestBodyJSON.URL)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url := fmt.Sprintf("%v/%v", s.cfg, shortenedURL)

	requestBodyJSON = ShortJSON{
		Result: url,
	}
	responseBody, err := json.Marshal(requestBodyJSON)
	if err != nil {
		log.Println("error marshaling JSON:", err)
		http.Error(w, "invalid data JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(responseBody))
}
