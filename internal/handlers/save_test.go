package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/gtngzlv/url-shortener/internal/errors"
	"github.com/gtngzlv/url-shortener/internal/models"
	"github.com/gtngzlv/url-shortener/internal/storage"
)

const testURL = "https://ya.ru"

func TestApp_PostURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name       string
		body       string
		methodType string
		dbErr      error
		want       want
	}{
		{
			name:       "POSTURL: First url save",
			methodType: http.MethodPost,
			body:       testURL,
			want: want{
				statusCode:  201,
				contentType: "plain/text",
			},
			dbErr: nil,
		},
		{
			name:       "POSTURL: empty body",
			methodType: http.MethodPost,
			body:       "",
			want: want{
				statusCode: 400,
			},
			dbErr: nil,
		},
		{
			name:       "POSTURL: already exist",
			methodType: http.MethodPost,
			body:       testURL,
			want: want{
				statusCode: 409,
			},
			dbErr: errors.ErrAlreadyExist,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			conf := returnTestConfig()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := storage.NewMockMyStorage(ctrl)
			dbMock.EXPECT().SaveFull(gomock.Any(), gomock.Any()).Return(models.URLInfo{}, tt.dbErr).AnyTimes()

			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.methodType, "/", body)
			w := httptest.NewRecorder()

			handler := NewApp(conf.router, conf.cfg, conf.log, dbMock)
			h := handler.PostURL

			h(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestApp_PostAPIShorten(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name        string
		body        string
		methodType  string
		contentType string
		dbErr       error
		want        want
	}{
		{
			name:        "POSTApiShorten: First url save",
			methodType:  http.MethodPost,
			contentType: "application/json",
			body:        `{"url":"123"}`,
			want: want{
				statusCode:  201,
				contentType: "application/json",
			},
			dbErr: nil,
		},
		{
			name:        "POSTApiShorten: conflict",
			methodType:  http.MethodPost,
			contentType: "application/json",
			body:        `{"url":"123"}`,
			want: want{
				statusCode:  409,
				contentType: "application/json",
			},
			dbErr: errors.ErrAlreadyExist,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			conf := returnTestConfig()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := storage.NewMockMyStorage(ctrl)
			dbMock.EXPECT().SaveFull(gomock.Any(), gomock.Any()).Return(models.URLInfo{}, tt.dbErr).AnyTimes()

			body := strings.NewReader(tt.body)
			request := httptest.NewRequest(tt.methodType, "/", body)
			w := httptest.NewRecorder()

			handler := NewApp(conf.router, conf.cfg, conf.log, dbMock)
			h := handler.PostAPIShorten

			h(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
