// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package util

import mock "github.com/stretchr/testify/mock"

// Hashable is an autogenerated mock type for the Hashable type
type Hashable struct {
	mock.Mock
}

// GetHash provides a mock function with given fields:
func (_m *Hashable) GetHash() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetHashBytes provides a mock function with given fields:
func (_m *Hashable) GetHashBytes() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}
