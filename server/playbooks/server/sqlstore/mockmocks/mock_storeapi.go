// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost/server/v8/playbooks/server/sqlstore (interfaces: StoreAPI)

// Package mock_sqlstore is a generated GoMock package.
package mock_sqlstore

import (
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStoreAPI is a mock of StoreAPI interface
type MockStoreAPI struct {
	ctrl     *gomock.Controller
	recorder *MockStoreAPIMockRecorder
}

// MockStoreAPIMockRecorder is the mock recorder for MockStoreAPI
type MockStoreAPIMockRecorder struct {
	mock *MockStoreAPI
}

// NewMockStoreAPI creates a new mock instance
func NewMockStoreAPI(ctrl *gomock.Controller) *MockStoreAPI {
	mock := &MockStoreAPI{ctrl: ctrl}
	mock.recorder = &MockStoreAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStoreAPI) EXPECT() *MockStoreAPIMockRecorder {
	return m.recorder
}

// DriverName mocks base method
func (m *MockStoreAPI) DriverName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DriverName")
	ret0, _ := ret[0].(string)
	return ret0
}

// DriverName indicates an expected call of DriverName
func (mr *MockStoreAPIMockRecorder) DriverName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DriverName", reflect.TypeOf((*MockStoreAPI)(nil).DriverName))
}

// GetMasterDB mocks base method
func (m *MockStoreAPI) GetMasterDB() (*sql.DB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMasterDB")
	ret0, _ := ret[0].(*sql.DB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMasterDB indicates an expected call of GetMasterDB
func (mr *MockStoreAPIMockRecorder) GetMasterDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMasterDB", reflect.TypeOf((*MockStoreAPI)(nil).GetMasterDB))
}
