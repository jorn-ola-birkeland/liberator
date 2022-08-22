// Code generated by mockery v2.12.1. DO NOT EDIT.

package nais_io_v1

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// MockAzureInterface is an autogenerated mock type for the AzureInterface type
type MockAzureInterface struct {
	mock.Mock
}

// GetApplication provides a mock function with given fields:
func (_m *MockAzureInterface) GetApplication() *AzureApplication {
	ret := _m.Called()

	var r0 *AzureApplication
	if rf, ok := ret.Get(0).(func() *AzureApplication); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AzureApplication)
		}
	}

	return r0
}

// GetSidecar provides a mock function with given fields:
func (_m *MockAzureInterface) GetSidecar() *AzureSidecar {
	ret := _m.Called()

	var r0 *AzureSidecar
	if rf, ok := ret.Get(0).(func() *AzureSidecar); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AzureSidecar)
		}
	}

	return r0
}

// NewMockAzureInterface creates a new instance of MockAzureInterface. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAzureInterface(t testing.TB) *MockAzureInterface {
	mock := &MockAzureInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}