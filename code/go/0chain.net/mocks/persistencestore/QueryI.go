// Code generated by mockery v2.6.0. DO NOT EDIT.

package persistencestore

import (
	context "context"

	gocql "github.com/gocql/gocql"

	mock "github.com/stretchr/testify/mock"

	persistencestore "0chain.net/core/persistencestore"
)

// QueryI is an autogenerated mock type for the QueryI type
type QueryI struct {
	mock.Mock
}

// AddAttempts provides a mock function with given fields: i, host
func (_m *QueryI) AddAttempts(i int, host *gocql.HostInfo) {
	_m.Called(i, host)
}

// AddLatency provides a mock function with given fields: l, host
func (_m *QueryI) AddLatency(l int64, host *gocql.HostInfo) {
	_m.Called(l, host)
}

// Attempts provides a mock function with given fields:
func (_m *QueryI) Attempts() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Bind provides a mock function with given fields: v
func (_m *QueryI) Bind(v ...interface{}) *gocql.Query {
	var _ca []interface{}
	_ca = append(_ca, v...)
	ret := _m.Called(_ca...)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(...interface{}) *gocql.Query); ok {
		r0 = rf(v...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// Cancel provides a mock function with given fields:
func (_m *QueryI) Cancel() {
	_m.Called()
}

// Consistency provides a mock function with given fields: c
func (_m *QueryI) Consistency(c gocql.Consistency) *gocql.Query {
	ret := _m.Called(c)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(gocql.Consistency) *gocql.Query); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// Context provides a mock function with given fields:
func (_m *QueryI) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// CustomPayload provides a mock function with given fields: customPayload
func (_m *QueryI) CustomPayload(customPayload map[string][]byte) *gocql.Query {
	ret := _m.Called(customPayload)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(map[string][]byte) *gocql.Query); ok {
		r0 = rf(customPayload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// DefaultTimestamp provides a mock function with given fields: enable
func (_m *QueryI) DefaultTimestamp(enable bool) *gocql.Query {
	ret := _m.Called(enable)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(bool) *gocql.Query); ok {
		r0 = rf(enable)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// Exec provides a mock function with given fields:
func (_m *QueryI) Exec() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetConsistency provides a mock function with given fields:
func (_m *QueryI) GetConsistency() gocql.Consistency {
	ret := _m.Called()

	var r0 gocql.Consistency
	if rf, ok := ret.Get(0).(func() gocql.Consistency); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(gocql.Consistency)
	}

	return r0
}

// GetRoutingKey provides a mock function with given fields:
func (_m *QueryI) GetRoutingKey() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Idempotent provides a mock function with given fields: value
func (_m *QueryI) Idempotent(value bool) *gocql.Query {
	ret := _m.Called(value)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(bool) *gocql.Query); ok {
		r0 = rf(value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// IsIdempotent provides a mock function with given fields:
func (_m *QueryI) IsIdempotent() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Iter provides a mock function with given fields:
func (_m *QueryI) Iter() persistencestore.IteratorI {
	ret := _m.Called()

	var r0 persistencestore.IteratorI
	if rf, ok := ret.Get(0).(func() persistencestore.IteratorI); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(persistencestore.IteratorI)
		}
	}

	return r0
}

// Keyspace provides a mock function with given fields:
func (_m *QueryI) Keyspace() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Latency provides a mock function with given fields:
func (_m *QueryI) Latency() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

// MapScan provides a mock function with given fields: m
func (_m *QueryI) MapScan(m map[string]interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MapScanCAS provides a mock function with given fields: dest
func (_m *QueryI) MapScanCAS(dest map[string]interface{}) (bool, error) {
	ret := _m.Called(dest)

	var r0 bool
	if rf, ok := ret.Get(0).(func(map[string]interface{}) bool); ok {
		r0 = rf(dest)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(map[string]interface{}) error); ok {
		r1 = rf(dest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NoSkipMetadata provides a mock function with given fields:
func (_m *QueryI) NoSkipMetadata() *gocql.Query {
	ret := _m.Called()

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func() *gocql.Query); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// Observer provides a mock function with given fields: observer
func (_m *QueryI) Observer(observer gocql.QueryObserver) *gocql.Query {
	ret := _m.Called(observer)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(gocql.QueryObserver) *gocql.Query); ok {
		r0 = rf(observer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// PageSize provides a mock function with given fields: n
func (_m *QueryI) PageSize(n int) *gocql.Query {
	ret := _m.Called(n)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(int) *gocql.Query); ok {
		r0 = rf(n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// PageState provides a mock function with given fields: state
func (_m *QueryI) PageState(state []byte) *gocql.Query {
	ret := _m.Called(state)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func([]byte) *gocql.Query); ok {
		r0 = rf(state)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// Prefetch provides a mock function with given fields: p
func (_m *QueryI) Prefetch(p float64) *gocql.Query {
	ret := _m.Called(p)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(float64) *gocql.Query); ok {
		r0 = rf(p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// Release provides a mock function with given fields:
func (_m *QueryI) Release() {
	_m.Called()
}

// RetryPolicy provides a mock function with given fields: r
func (_m *QueryI) RetryPolicy(r gocql.RetryPolicy) *gocql.Query {
	ret := _m.Called(r)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(gocql.RetryPolicy) *gocql.Query); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// RoutingKey provides a mock function with given fields: routingKey
func (_m *QueryI) RoutingKey(routingKey []byte) *gocql.Query {
	ret := _m.Called(routingKey)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func([]byte) *gocql.Query); ok {
		r0 = rf(routingKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// Scan provides a mock function with given fields: dest
func (_m *QueryI) Scan(dest ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, dest...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(dest...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ScanCAS provides a mock function with given fields: dest
func (_m *QueryI) ScanCAS(dest ...interface{}) (bool, error) {
	var _ca []interface{}
	_ca = append(_ca, dest...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(...interface{}) bool); ok {
		r0 = rf(dest...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...interface{}) error); ok {
		r1 = rf(dest...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SerialConsistency provides a mock function with given fields: cons
func (_m *QueryI) SerialConsistency(cons gocql.SerialConsistency) *gocql.Query {
	ret := _m.Called(cons)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(gocql.SerialConsistency) *gocql.Query); ok {
		r0 = rf(cons)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// SetConsistency provides a mock function with given fields: c
func (_m *QueryI) SetConsistency(c gocql.Consistency) {
	_m.Called(c)
}

// SetSpeculativeExecutionPolicy provides a mock function with given fields: sp
func (_m *QueryI) SetSpeculativeExecutionPolicy(sp gocql.SpeculativeExecutionPolicy) *gocql.Query {
	ret := _m.Called(sp)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(gocql.SpeculativeExecutionPolicy) *gocql.Query); ok {
		r0 = rf(sp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// Statement provides a mock function with given fields:
func (_m *QueryI) Statement() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *QueryI) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Trace provides a mock function with given fields: trace
func (_m *QueryI) Trace(trace gocql.Tracer) *gocql.Query {
	ret := _m.Called(trace)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(gocql.Tracer) *gocql.Query); ok {
		r0 = rf(trace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// WithContext provides a mock function with given fields: ctx
func (_m *QueryI) WithContext(ctx context.Context) *gocql.Query {
	ret := _m.Called(ctx)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(context.Context) *gocql.Query); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}

// WithTimestamp provides a mock function with given fields: timestamp
func (_m *QueryI) WithTimestamp(timestamp int64) *gocql.Query {
	ret := _m.Called(timestamp)

	var r0 *gocql.Query
	if rf, ok := ret.Get(0).(func(int64) *gocql.Query); ok {
		r0 = rf(timestamp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gocql.Query)
		}
	}

	return r0
}
