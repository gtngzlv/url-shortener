package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/gtngzlv/url-shortener/internal/storage"
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
			cfg.DatabaseDSN = "postgresql://newuser:password@localhost/postgres?sslmode=disable"

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dbMock := storage.NewMockMyStorage(ctrl)

			handler := NewApp(router, cfg, log, dbMock)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.URL.Path = tt.query
			w := httptest.NewRecorder()
			dbMock.EXPECT().GetByShort("").Return("ya.ru", nil)
			handler.GetURL(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedCode)
		})
	}
}
