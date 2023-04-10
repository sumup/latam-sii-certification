// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IAdapter is an autogenerated mock type for the IAdapter type
type IAdapter struct {
	mock.Mock
}

// Post provides a mock function with given fields: ctx, url, payload
func (_m *IAdapter) Post(ctx context.Context, url string, headers map[string]string, payload []byte) ([]byte, error) {
	ret := _m.Called(ctx, url, payload)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) []byte); ok {
		r0 = rf(ctx, url, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []byte) error); ok {
		r1 = rf(ctx, url, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewIAdapter creates a new instance of IAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIAdapter(t mockConstructorTestingTNewIAdapter) *IAdapter {
	mock := &IAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
