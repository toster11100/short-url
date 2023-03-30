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
		name            string
		target          string
		method          string
		requestBody     string
		wantBody        string
		wantContentType string
		wantCode        int
		location        string
	}{
		{
			name:            "case ok 1: post",
			target:          "/",
			method:          http.MethodPost,
			requestBody:     "https://www.google.com",
			wantBody:        "http://localhost:8080/1",
			wantContentType: "text/html",
			wantCode:        http.StatusCreated,
			location:        "",
		},
		{
			name:            "case ok 2: post",
			target:          "/",
			method:          http.MethodPost,
			requestBody:     "https://ya.ru",
			wantBody:        "http://localhost:8080/2",
			wantContentType: "text/html",
			wantCode:        http.StatusCreated,
			location:        "",
		},
		{
			name:            "case ok 3: get",
			target:          "/1",
			method:          http.MethodGet,
			requestBody:     "",
			wantBody:        "",
			wantContentType: "text/html",
			wantCode:        http.StatusTemporaryRedirect,
			location:        "https://www.google.com",
		},
		{
			name:            "case ok 4: api/shortned",
			target:          "/api/shorten",
			method:          http.MethodPost,
			requestBody:     "{\"url\":\"http://ya.ru\"}",
			wantBody:        "{\"result\":\"http://localhost:8080/3\"}",
			wantContentType: "application/json",
			wantCode:        http.StatusCreated,
			location:        "",
		},
		{
			name:            "case don't ok 1: wrong url",
			target:          "/",
			method:          http.MethodPost,
			requestBody:     "123",
			wantBody:        "this \"123\" is not URL\n",
			wantContentType: "text/plain; charset=utf-8",
			wantCode:        http.StatusBadRequest,
			location:        "",
		},
		{
			name:            "case don't ok 2: wrong method",
			target:          "/",
			method:          http.MethodPut,
			requestBody:     "",
			wantBody:        "",
			wantContentType: "",
			wantCode:        http.StatusMethodNotAllowed,
			location:        "",
		},
		{
			name:            "case don't ok 3: wrong id get",
			target:          "/9223372036854775808",
			method:          http.MethodGet,
			requestBody:     "",
			wantBody:        "this ID: 9223372036854775808 is not valid\n",
			wantContentType: "text/plain; charset=utf-8",
			wantCode:        http.StatusBadRequest,
			location:        "",
		},
		{
			name:            "case don't ok 4: wrong id",
			target:          "/10",
			method:          http.MethodGet,
			requestBody:     "",
			wantBody:        "URL with ID 10 is not found in URLMap\n",
			wantContentType: "text/plain; charset=utf-8",
			wantCode:        http.StatusNotFound,
			location:        "",
		},
		{
			name:            "case don't ok 5: invalid JSON",
			target:          "/api/shorten",
			method:          http.MethodPost,
			requestBody:     "http://ya.ru",
			wantBody:        "invalid data JSON\n",
			wantContentType: "text/plain; charset=utf-8",
			wantCode:        http.StatusBadRequest,
			location:        "",
		},
		{
			name:            "case don't ok 6: empty URL",
			target:          "/api/shorten",
			method:          http.MethodPost,
			requestBody:     "{\"url\":\"\"}",
			wantBody:        "URL field is required\n",
			wantContentType: "text/plain; charset=utf-8",
			wantCode:        http.StatusBadRequest,
			location:        "",
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
			assert.Equal(t, tt.wantContentType, response.Header.Get("Content-Type"))
		})
	}
}
