// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/storage.go

// Package storage is a generated GoMock package.
package storage

import (
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/gtngzlv/url-shortener/internal/models"
)

// MockMyStorage is a mock of MyStorage interface.
type MockMyStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMyStorageMockRecorder
}

// MockMyStorageMockRecorder is the mock recorder for MockMyStorage.
type MockMyStorageMockRecorder struct {
	mock *MockMyStorage
}

// NewMockMyStorage creates a new mock instance.
func NewMockMyStorage(ctrl *gomock.Controller) *MockMyStorage {
	mock := &MockMyStorage{ctrl: ctrl}
	mock.recorder = &MockMyStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMyStorage) EXPECT() *MockMyStorageMockRecorder {
	return m.recorder
}

// Batch mocks base method.
func (m *MockMyStorage) Batch(userID string, entities []models.URLInfo) ([]models.URLInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Batch", userID, entities)
	ret0, _ := ret[0].([]models.URLInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Batch indicates an expected call of Batch.
func (mr *MockMyStorageMockRecorder) Batch(userID, entities interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Batch", reflect.TypeOf((*MockMyStorage)(nil).Batch), userID, entities)
}

// DeleteByUserIDAndShort mocks base method.
func (m *MockMyStorage) DeleteByUserIDAndShort(userID, shortURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUserIDAndShort", userID, shortURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUserIDAndShort indicates an expected call of DeleteByUserIDAndShort.
func (mr *MockMyStorageMockRecorder) DeleteByUserIDAndShort(userID, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUserIDAndShort", reflect.TypeOf((*MockMyStorage)(nil).DeleteByUserIDAndShort), userID, shortURL)
}

// GetBatchByUserID mocks base method.
func (m *MockMyStorage) GetBatchByUserID(userID string) ([]models.URLInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatchByUserID", userID)
	ret0, _ := ret[0].([]models.URLInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatchByUserID indicates an expected call of GetBatchByUserID.
func (mr *MockMyStorageMockRecorder) GetBatchByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatchByUserID", reflect.TypeOf((*MockMyStorage)(nil).GetBatchByUserID), userID)
}

// GetByShort mocks base method.
func (m *MockMyStorage) GetByShort(shortURL string) (models.URLInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByShort", shortURL)
	ret0, _ := ret[0].(models.URLInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByShort indicates an expected call of GetByShort.
func (mr *MockMyStorageMockRecorder) GetByShort(shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByShort", reflect.TypeOf((*MockMyStorage)(nil).GetByShort), shortURL)
}

// GetStatistic mocks base method.
func (m *MockMyStorage) GetStatistic() *models.Statistic {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatistic")
	ret0, _ := ret[0].(*models.Statistic)
	return ret0
}

// GetStatistic indicates an expected call of GetStatistic.
func (mr *MockMyStorageMockRecorder) GetStatistic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatistic", reflect.TypeOf((*MockMyStorage)(nil).GetStatistic))
}

// Ping mocks base method.
func (m *MockMyStorage) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockMyStorageMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockMyStorage)(nil).Ping))
}

// SaveFull mocks base method.
func (m *MockMyStorage) SaveFull(userID, fullURL string) (models.URLInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFull", userID, fullURL)
	ret0, _ := ret[0].(models.URLInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveFull indicates an expected call of SaveFull.
func (mr *MockMyStorageMockRecorder) SaveFull(userID, fullURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFull", reflect.TypeOf((*MockMyStorage)(nil).SaveFull), userID, fullURL)
}