// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make plugin-mocks`.

package plugintest

import (
	model "github.com/mattermost/mattermost-server/model"
	plugin "github.com/mattermost/mattermost-server/plugin"
	mock "github.com/stretchr/testify/mock"
)

// Helpers is an autogenerated mock type for the Helpers type
type Helpers struct {
	mock.Mock
}

// EnsureBot provides a mock function with given fields: bot, options
func (_m *Helpers) EnsureBot(bot *model.Bot, options ...plugin.EnsureBotOption) (string, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, bot)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	if rf, ok := ret.Get(0).(func(*model.Bot, ...plugin.EnsureBotOption) string); ok {
		r0 = rf(bot, options...)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Bot, ...plugin.EnsureBotOption) error); ok {
		r1 = rf(bot, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InstallPluginFromURL provides a mock function with given fields: url, replace
func (_m *Helpers) InstallPluginFromURL(url string, replace bool) (*model.Manifest, error) {
	ret := _m.Called(url, replace)

	var r0 *model.Manifest
	if rf, ok := ret.Get(0).(func(string, bool) *model.Manifest); ok {
		r0 = rf(url, replace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Manifest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(url, replace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KVCompareAndDeleteJSON provides a mock function with given fields: key, oldValue
func (_m *Helpers) KVCompareAndDeleteJSON(key string, oldValue interface{}) (bool, error) {
	ret := _m.Called(key, oldValue)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, interface{}) bool); ok {
		r0 = rf(key, oldValue)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(key, oldValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KVCompareAndSetJSON provides a mock function with given fields: key, oldValue, newValue
func (_m *Helpers) KVCompareAndSetJSON(key string, oldValue interface{}, newValue interface{}) (bool, error) {
	ret := _m.Called(key, oldValue, newValue)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) bool); ok {
		r0 = rf(key, oldValue, newValue)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, interface{}, interface{}) error); ok {
		r1 = rf(key, oldValue, newValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KVGetJSON provides a mock function with given fields: key, value
func (_m *Helpers) KVGetJSON(key string, value interface{}) (bool, error) {
	ret := _m.Called(key, value)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, interface{}) bool); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(key, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KVSetJSON provides a mock function with given fields: key, value
func (_m *Helpers) KVSetJSON(key string, value interface{}) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// KVSetWithExpiryJSON provides a mock function with given fields: key, value, expireInSeconds
func (_m *Helpers) KVSetWithExpiryJSON(key string, value interface{}, expireInSeconds int64) error {
	ret := _m.Called(key, value, expireInSeconds)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, int64) error); ok {
		r0 = rf(key, value, expireInSeconds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ShouldProcessMessage provides a mock function with given fields: post, options
func (_m *Helpers) ShouldProcessMessage(post *model.Post, options ...plugin.ShouldProcessMessageOption) (bool, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, post)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*model.Post, ...plugin.ShouldProcessMessageOption) bool); ok {
		r0 = rf(post, options...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Post, ...plugin.ShouldProcessMessageOption) error); ok {
		r1 = rf(post, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
