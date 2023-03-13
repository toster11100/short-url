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
	srv := &Server{
		urlMap:  make(storage.Rep),
		id:      1,
		handler: nil,
	}
	tests := []struct {
		name        string
		requestBody string
		wantCode    int
		wantBody    string
		method      string
		target      string
		location    string
	}{
		{
			name:        "case pos 1",
			requestBody: "https://www.google.com",
			wantCode:    http.StatusCreated,
			wantBody:    "http://localhost:8080/1",
			method:      http.MethodPost,
			target:      "/",
			location:    "",
		},
		{
			name:        "case pos 2",
			requestBody: "https://ya.ru",
			wantCode:    http.StatusCreated,
			wantBody:    "http://localhost:8080/2",
			method:      http.MethodPost,
			target:      "/",
			location:    "",
		},
		{
			name:        "case pos 3",
			requestBody: "",
			wantCode:    http.StatusTemporaryRedirect,
			wantBody:    "",
			method:      http.MethodGet,
			target:      "/1",
			location:    "https://www.google.com",
		},
		{
			name:        "case neg 1",
			requestBody: "123",
			wantCode:    http.StatusBadRequest,
			wantBody:    "this is not URL\n",
			method:      http.MethodPost,
			target:      "/",
			location:    "",
		},
		{
			name:        "case neg 2",
			requestBody: "",
			wantCode:    http.StatusBadRequest,
			wantBody:    "invalid method\n",
			method:      http.MethodPut,
			target:      "/",
			location:    "",
		},
		{
			name:        "case neg 3",
			requestBody: "",
			wantCode:    http.StatusBadRequest,
			wantBody:    "invalid short id\n",
			method:      http.MethodGet,
			target:      "/0",
			location:    "",
		},
		{
			name:        "case neg 4",
			requestBody: "",
			wantCode:    http.StatusBadRequest,
			wantBody:    "invalid short id\n",
			method:      http.MethodGet,
			target:      "/10",
			location:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(srv.CheckMethod)
			h(w, request)
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
