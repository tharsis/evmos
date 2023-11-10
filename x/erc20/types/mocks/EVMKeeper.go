// Code generated by mockery v2.36.1. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	core "github.com/ethereum/go-ethereum/core"

	evmtypes "github.com/evmos/evmos/v15/x/evm/types"

	mock "github.com/stretchr/testify/mock"

	statedb "github.com/evmos/evmos/v15/x/evm/statedb"

	types "github.com/cosmos/cosmos-sdk/types"

	vm "github.com/ethereum/go-ethereum/core/vm"
)

// EVMKeeper is an autogenerated mock type for the EVMKeeper type
type EVMKeeper struct {
	mock.Mock
}

// AddEVMExtensions provides a mock function with given fields: ctx, precompiles
func (_m *EVMKeeper) AddEVMExtensions(ctx types.Context, precompiles ...vm.PrecompiledContract) error {
	_va := make([]interface{}, len(precompiles))
	for _i := range precompiles {
		_va[_i] = precompiles[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, ...vm.PrecompiledContract) error); ok {
		r0 = rf(ctx, precompiles...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ApplyMessage provides a mock function with given fields: ctx, msg, tracer, commit
func (_m *EVMKeeper) ApplyMessage(ctx types.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtypes.MsgEthereumTxResponse, error) {
	ret := _m.Called(ctx, msg, tracer, commit)

	var r0 *evmtypes.MsgEthereumTxResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, core.Message, vm.EVMLogger, bool) (*evmtypes.MsgEthereumTxResponse, error)); ok {
		return rf(ctx, msg, tracer, commit)
	}
	if rf, ok := ret.Get(0).(func(types.Context, core.Message, vm.EVMLogger, bool) *evmtypes.MsgEthereumTxResponse); ok {
		r0 = rf(ctx, msg, tracer, commit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*evmtypes.MsgEthereumTxResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, core.Message, vm.EVMLogger, bool) error); ok {
		r1 = rf(ctx, msg, tracer, commit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccount provides a mock function with given fields: ctx, addr
func (_m *EVMKeeper) DeleteAccount(ctx types.Context, addr common.Address) error {
	ret := _m.Called(ctx, addr)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address) error); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EstimateGasInternal provides a mock function with given fields: c, req, fromType
func (_m *EVMKeeper) EstimateGasInternal(c context.Context, req *evmtypes.EthCallRequest, fromType evmtypes.CallType) (*evmtypes.EstimateGasResponse, error) {
	ret := _m.Called(c, req, fromType)

	var r0 *evmtypes.EstimateGasResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *evmtypes.EthCallRequest, evmtypes.CallType) (*evmtypes.EstimateGasResponse, error)); ok {
		return rf(c, req, fromType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *evmtypes.EthCallRequest, evmtypes.CallType) *evmtypes.EstimateGasResponse); ok {
		r0 = rf(c, req, fromType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*evmtypes.EstimateGasResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *evmtypes.EthCallRequest, evmtypes.CallType) error); ok {
		r1 = rf(c, req, fromType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountWithoutBalance provides a mock function with given fields: ctx, addr
func (_m *EVMKeeper) GetAccountWithoutBalance(ctx types.Context, addr common.Address) *statedb.Account {
	ret := _m.Called(ctx, addr)

	var r0 *statedb.Account
	if rf, ok := ret.Get(0).(func(types.Context, common.Address) *statedb.Account); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*statedb.Account)
		}
	}

	return r0
}

// GetParams provides a mock function with given fields: ctx
func (_m *EVMKeeper) GetParams(ctx types.Context) evmtypes.Params {
	ret := _m.Called(ctx)

	var r0 evmtypes.Params
	if rf, ok := ret.Get(0).(func(types.Context) evmtypes.Params); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(evmtypes.Params)
	}

	return r0
}

// NewEVMKeeper creates a new instance of EVMKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEVMKeeper(t interface {
	mock.TestingT
	Cleanup(func())
},
) *EVMKeeper {
	mock := &EVMKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
