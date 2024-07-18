// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/evmos/evmos/v18/x/erc20/types"
	mock "github.com/stretchr/testify/mock"
)

// QueryServer is an autogenerated mock type for the QueryServer type
type QueryServer struct {
	mock.Mock
}

// Params provides a mock function with given fields: _a0, _a1
func (_m *QueryServer) Params(_a0 context.Context, _a1 *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Params")
	}

	var r0 *types.QueryParamsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryParamsRequest) (*types.QueryParamsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryParamsRequest) *types.QueryParamsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryParamsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryParamsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenPair provides a mock function with given fields: _a0, _a1
func (_m *QueryServer) TokenPair(_a0 context.Context, _a1 *types.QueryTokenPairRequest) (*types.QueryTokenPairResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for TokenPair")
	}

	var r0 *types.QueryTokenPairResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryTokenPairRequest) (*types.QueryTokenPairResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryTokenPairRequest) *types.QueryTokenPairResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryTokenPairResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryTokenPairRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenPairs provides a mock function with given fields: _a0, _a1
func (_m *QueryServer) TokenPairs(_a0 context.Context, _a1 *types.QueryTokenPairsRequest) (*types.QueryTokenPairsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for TokenPairs")
	}

	var r0 *types.QueryTokenPairsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryTokenPairsRequest) (*types.QueryTokenPairsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryTokenPairsRequest) *types.QueryTokenPairsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryTokenPairsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryTokenPairsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewQueryServer creates a new instance of QueryServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQueryServer(t interface {
	mock.TestingT
	Cleanup(func())
},
) *QueryServer {
	mock := &QueryServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
