// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entities "github.com/sumup/sii-certification/internal/entities"
)

// IGateway is an autogenerated mock type for the IGateway type
type IGateway struct {
	mock.Mock
}

// GetSeed provides a mock function with given fields: ctx
func (_m *IGateway) GetSeed(ctx context.Context) (string, string, error) {
	ret := _m.Called(ctx)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context) string); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetToken provides a mock function with given fields: ctx, seed
func (_m *IGateway) GetToken(ctx context.Context, seed string) (string, string, error) {
	ret := _m.Called(ctx, seed)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, seed)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, seed)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, seed)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SendMany provides a mock function with given fields: ctx, batch
func (_m *IGateway) SendMany(ctx context.Context, batch []entities.Batch) (bool, string, []entities.Batch, error) {
	ret := _m.Called(ctx, batch)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, []entities.Batch) bool); ok {
		r0 = rf(ctx, batch)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, []entities.Batch) string); ok {
		r1 = rf(ctx, batch)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 []entities.Batch
	if rf, ok := ret.Get(2).(func(context.Context, []entities.Batch) []entities.Batch); ok {
		r2 = rf(ctx, batch)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).([]entities.Batch)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, []entities.Batch) error); ok {
		r3 = rf(ctx, batch)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

type mockConstructorTestingTNewIGateway interface {
	mock.TestingT
	Cleanup(func())
}

// NewIGateway creates a new instance of IGateway. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIGateway(t mockConstructorTestingTNewIGateway) *IGateway {
	mock := &IGateway{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
