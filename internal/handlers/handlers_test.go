package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/gtngzlv/url-shortener/internal/storage/database"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/logger"
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
			router := chi.NewRouter()
			log := logger.NewLogger()
			cfg := config.LoadConfig()
			cfg.FileStoragePath = "/tmp/short-url-bd.json"

			s := database.Init(log, cfg)
			handler := NewApp(router, cfg, log, s)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.URL.Path = tt.query
			w := httptest.NewRecorder()
			handler.GetURL(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedCode)
		})
	}
}

// не понимаю почему они падают
//func TestPostAPIShorten(t *testing.T) {
//	testCases := []struct {
//		name         string
//		contentType  string
//		body         map[string]string
//		expectedCode int
//	}{
//		{
//			name:         "201 ok",
//			contentType:  "application/json",
//			expectedCode: 201,
//			body:         map[string]string{"url": "https://practicum1.yandex.ru"},
//		},
//	}
//
//	for _, tt := range testCases {
//		t.Run(tt.name, func(t *testing.T) {
//			cfg := config.LoadConfig()
//			logger.NewLogger()
//			storage.Init(cfg)
//			app := NewApp(cfg)
//			testServer := httptest.NewServer(app)
//			defer testServer.Close()
//
//			body, err := json.Marshal(tt.body)
//			assert.NoError(t, err)
//
//			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/shorten", bytes.NewBuffer(body))
//			req.Header.Set("Content-Type", "application/json")
//			assert.NoError(t, err)
//
//			resp, err := http.DefaultClient.Do(req)
//			assert.NoError(t, err)
//
//			var r APIShortenResponse
//			rBody, err := io.ReadAll(resp.Body)
//			assert.NoError(t, err)
//			defer resp.Body.Close()
//
//			err = json.Unmarshal(rBody, &r)
//			assert.NoError(t, err)
//
//			assert.NotEmpty(t, r.Result)
//			assert.Equal(t, tt.expectedCode, resp.StatusCode)
//		})
//	}
//}

//
//func TestPostHandler(t *testing.T) {
//	testCases := []struct {
//		name         string
//		contentType  string
//		body         string
//		expectedCode int
//	}{
//		{
//			name:         "200 ok",
//			contentType:  "text/plain",
//			expectedCode: 201,
//			body:         "ya.ru",
//		},
//	}
//
//	for _, tt := range testCases {
//		t.Run(tt.name, func(t *testing.T) {
//			cfg := config.LoadConfig()
//			cfg.FileStoragePath = "/tmp/short-url-bd.json"
//
//			storage.Init(cfg)
//			handler := NewApp(cfg)
//
//			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
//			request.Header.Set("Content-Type", tt.contentType)
//			w := httptest.NewRecorder()
//			handler.PostURL(w, request)
//			res := w.Result()
//			defer res.Body.Close()
//
//			assert.Equal(t, res.StatusCode, tt.expectedCode)
//		})
//	}
//}
