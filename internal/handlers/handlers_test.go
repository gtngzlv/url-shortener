package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func TestPostHandler(t *testing.T) {
	testCases := []struct {
		name                string
		contentType         string
		body                string
		expectedCode        int
		expectedBodyIsEmpty bool
	}{
		{
			name:         "200 ok",
			contentType:  "text/plain",
			expectedCode: 201,
			body:         "ya.ru",
		},
		{
			name:         "400",
			contentType:  "applicaiton/json",
			expectedCode: 400,
			body:         "google.com",
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
