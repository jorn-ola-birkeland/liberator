// Code generated by mockery v2.12.1. DO NOT EDIT.

package service

import (
	aiven "github.com/aiven/aiven-go-client"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// MockInterface is an autogenerated mock type for the Interface type
type MockInterface struct {
	mock.Mock
}

// List provides a mock function with given fields: project
func (_m *MockInterface) List(project string) ([]*aiven.Service, error) {
	ret := _m.Called(project)

	var r0 []*aiven.Service
	if rf, ok := ret.Get(0).(func(string) []*aiven.Service); ok {
		r0 = rf(project)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aiven.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(project)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockInterface creates a new instance of MockInterface. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockInterface(t testing.TB) *MockInterface {
	mock := &MockInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}