package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	testCases := []struct {
		name         string
		query        string
		expectedCode int
	}{
		{
			name:         "307 ok",
			query:        "12345",
			expectedCode: 307,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.URL.Path = tt.query
			w := httptest.NewRecorder()
			GetURL(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedCode)
		})
	}
}

func TestPostAPIShorten(t *testing.T) {
	testCases := []struct {
		name         string
		contentType  string
		body         map[string]string
		expectedCode int
	}{
		{
			name:         "201 ok",
			contentType:  "application/json",
			expectedCode: 201,
			body:         map[string]string{"url": "https://practicum.yandex.ru"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			PostAPIShorten(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedCode)
		})
	}
}

func TestPostHandler(t *testing.T) {
	testCases := []struct {
		name         string
		contentType  string
		body         string
		expectedCode int
	}{
		{
			name:         "200 ok",
			contentType:  "text/plain",
			expectedCode: 201,
			body:         "ya.ru",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			PostURL(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedCode)
		})
	}
}
