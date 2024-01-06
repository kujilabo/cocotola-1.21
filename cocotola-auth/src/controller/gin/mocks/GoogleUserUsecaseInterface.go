// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"

	mock "github.com/stretchr/testify/mock"
)

// GoogleUserUsecaseInterface is an autogenerated mock type for the GoogleUserUsecaseInterface type
type GoogleUserUsecaseInterface struct {
	mock.Mock
}

type GoogleUserUsecaseInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *GoogleUserUsecaseInterface) EXPECT() *GoogleUserUsecaseInterface_Expecter {
	return &GoogleUserUsecaseInterface_Expecter{mock: &_m.Mock}
}

// Authorize provides a mock function with given fields: ctx, code, organizationName
func (_m *GoogleUserUsecaseInterface) Authorize(ctx context.Context, code string, organizationName string) (*domain.AuthTokenSet, error) {
	ret := _m.Called(ctx, code, organizationName)

	var r0 *domain.AuthTokenSet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*domain.AuthTokenSet, error)); ok {
		return rf(ctx, code, organizationName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *domain.AuthTokenSet); ok {
		r0 = rf(ctx, code, organizationName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.AuthTokenSet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, code, organizationName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GoogleUserUsecaseInterface_Authorize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authorize'
type GoogleUserUsecaseInterface_Authorize_Call struct {
	*mock.Call
}

// Authorize is a helper method to define mock.On call
//   - ctx context.Context
//   - code string
//   - organizationName string
func (_e *GoogleUserUsecaseInterface_Expecter) Authorize(ctx interface{}, code interface{}, organizationName interface{}) *GoogleUserUsecaseInterface_Authorize_Call {
	return &GoogleUserUsecaseInterface_Authorize_Call{Call: _e.mock.On("Authorize", ctx, code, organizationName)}
}

func (_c *GoogleUserUsecaseInterface_Authorize_Call) Run(run func(ctx context.Context, code string, organizationName string)) *GoogleUserUsecaseInterface_Authorize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *GoogleUserUsecaseInterface_Authorize_Call) Return(_a0 *domain.AuthTokenSet, _a1 error) *GoogleUserUsecaseInterface_Authorize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GoogleUserUsecaseInterface_Authorize_Call) RunAndReturn(run func(context.Context, string, string) (*domain.AuthTokenSet, error)) *GoogleUserUsecaseInterface_Authorize_Call {
	_c.Call.Return(run)
	return _c
}

// NewGoogleUserUsecaseInterface creates a new instance of GoogleUserUsecaseInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGoogleUserUsecaseInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *GoogleUserUsecaseInterface {
	mock := &GoogleUserUsecaseInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}