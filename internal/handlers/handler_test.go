package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"github.com/toster11100/shortUrl.git/internal/storage"
)

func TestServer_createShortURL(t *testing.T) {
	repo := storage.New()
	srv := New(repo)
	tests := []struct {
		name        string
		target      string
		wantCode    int
		requestBody string
		wantBody    string
		method      string
		location    string
	}{
		{
			name:        "case pos 1",
			target:      "/",
			wantCode:    http.StatusCreated,
			requestBody: "https://www.google.com",
			wantBody:    "http://localhost:8080/1",
			method:      http.MethodPost,
			location:    "",
		},
		{
			name:        "case pos 2",
			target:      "/",
			wantCode:    http.StatusCreated,
			requestBody: "https://ya.ru",
			wantBody:    "http://localhost:8080/2",
			method:      http.MethodPost,
			location:    "",
		},
		{
			name:        "case pos 3",
			target:      "/1",
			wantCode:    http.StatusTemporaryRedirect,
			requestBody: "",
			wantBody:    "",
			method:      http.MethodGet,
			location:    "https://www.google.com",
		},
		{
			name:        "case neg 1",
			target:      "/",
			wantCode:    http.StatusBadRequest,
			requestBody: "123",
			wantBody:    "this \"123\" is not URL\n",
			method:      http.MethodPost,
			location:    "",
		},
		{
			name:        "case neg 2",
			target:      "/",
			wantCode:    http.StatusMethodNotAllowed,
			requestBody: "",
			wantBody:    "",
			method:      http.MethodPut,
			location:    "",
		},
		{
			name:        "case neg 3",
			target:      "/9223372036854775808",
			wantCode:    http.StatusBadRequest,
			requestBody: "",
			wantBody:    "this ID: 9223372036854775808 is not valid\n",
			method:      http.MethodGet,
			location:    "",
		},
		{
			name:        "case neg 4",
			target:      "/10",
			wantCode:    http.StatusNotFound,
			requestBody: "",
			wantBody:    "URL with ID 10 is not found in URLMap\n",
			method:      http.MethodGet,
			location:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, request)
			response := w.Result()
			assert.Equal(t, tt.wantCode, response.StatusCode)
			userResult, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			defer response.Body.Close()
			assert.Equal(t, tt.wantBody, string(userResult))
			assert.Equal(t, tt.location, response.Header.Get("Location"))
		})
	}
}
