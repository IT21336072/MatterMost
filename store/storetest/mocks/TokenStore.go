// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/mattermost/mattermost-server/model"
import store "github.com/mattermost/mattermost-server/store"

// TokenStore is an autogenerated mock type for the TokenStore type
type TokenStore struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields:
func (_m *TokenStore) Cleanup() {
	_m.Called()
}

// Delete provides a mock function with given fields: token
func (_m *TokenStore) Delete(token string) *model.AppError {
	ret := _m.Called(token)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// GetByToken provides a mock function with given fields: token
func (_m *TokenStore) GetByToken(token string) (*model.Token, *model.AppError) {
	ret := _m.Called(token)

	var r0 *model.Token
	if rf, ok := ret.Get(0).(func(string) *model.Token); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Token)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// RemoveAllTokensByType provides a mock function with given fields: tokenType
func (_m *TokenStore) RemoveAllTokensByType(tokenType string) store.StoreChannel {
	ret := _m.Called(tokenType)

	var r0 store.StoreChannel
	if rf, ok := ret.Get(0).(func(string) store.StoreChannel); ok {
		r0 = rf(tokenType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.StoreChannel)
		}
	}

	return r0
}

// Save provides a mock function with given fields: recovery
func (_m *TokenStore) Save(recovery *model.Token) *model.AppError {
	ret := _m.Called(recovery)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.Token) *model.AppError); ok {
		r0 = rf(recovery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}
