package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.uber.org/zap"

	err "github.com/gtngzlv/url-shortener/internal/errors"
	"github.com/gtngzlv/url-shortener/internal/models"
	"github.com/gtngzlv/url-shortener/internal/storage"

	"github.com/stretchr/testify/assert"

	"github.com/gtngzlv/url-shortener/internal/config"
	"github.com/gtngzlv/url-shortener/internal/logger"
)

type TestConfig struct {
	router *chi.Mux
	log    zap.SugaredLogger
	cfg    *config.AppConfig
}

func returnTestConfig() TestConfig {
	return TestConfig{
		router: chi.NewRouter(),
		log:    logger.NewLogger(),
		cfg: &config.AppConfig{
			FileStoragePath: "/1/1.json",
			DatabaseDSN:     "pg://",
		},
	}
}

func TestGetHandlerSuccess(t *testing.T) {
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
			conf := returnTestConfig()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dbMock := storage.NewMockMyStorage(ctrl)

			handler := NewApp(conf.router, conf.cfg, conf.log, dbMock)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.URL.Path = tt.query
			w := httptest.NewRecorder()

			mockedDBResult := models.URLInfo{
				OriginalURL: "ya.ru",
			}
			dbMock.EXPECT().GetByShort("").Return(mockedDBResult, nil)
			handler.GetURL(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedCode)
		})
	}
}

func TestGetHandlerError(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		query    string
		dbErr    error
		wantCode int
		dbResp   models.URLInfo
	}{{
		name:     "empty db info",
		query:    "12345",
		dbErr:    errors.New("error"),
		wantCode: http.StatusBadRequest,
		dbResp:   models.URLInfo{},
	},
		{
			name:     "url is deleted",
			query:    "12345",
			dbErr:    nil,
			wantCode: http.StatusGone,
			dbResp: models.URLInfo{
				OriginalURL: "ya.ru",
				IsDeleted:   1,
			},
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			conf := returnTestConfig()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dbMock := storage.NewMockMyStorage(ctrl)

			handler := NewApp(conf.router, conf.cfg, conf.log, dbMock)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.URL.Path = tt.query
			w := httptest.NewRecorder()

			dbMock.EXPECT().GetByShort("").Return(tt.dbResp, tt.dbErr)

			handler.GetURL(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantCode)
		})
	}
}

func TestApp_GetURLsDBError(t *testing.T) {
	conf := returnTestConfig()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbMock := storage.NewMockMyStorage(ctrl)

	handler := NewApp(conf.router, conf.cfg, conf.log, dbMock)
	request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	w := httptest.NewRecorder()

	dbMock.EXPECT().GetBatchByUserID(gomock.Any()).Return([]models.URLInfo{}, err.ErrNoBatchByUserID)
	handler.GetURLs(w, request)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusNoContent)
}

func TestApp_GetURLsDBSuccess(t *testing.T) {
	conf := returnTestConfig()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbMock := storage.NewMockMyStorage(ctrl)

	handler := NewApp(conf.router, conf.cfg, conf.log, dbMock)
	request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	w := httptest.NewRecorder()

	dbMock.EXPECT().GetBatchByUserID(gomock.Any()).Return([]models.URLInfo{
		{
			UUID:        uuid.NewString(),
			OriginalURL: testURL,
		},
	}, nil)
	handler.GetURLs(w, request)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
