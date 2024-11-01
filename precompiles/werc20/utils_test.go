package werc20_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/evmos/evmos/v20/testutil/integration/evmos/factory"
	"github.com/evmos/evmos/v20/testutil/integration/evmos/keyring"
	evmtypes "github.com/evmos/evmos/v20/x/evm/types"
)

// callType constants to differentiate between
// the different types of call to the precompile.
type callType int

const (
	directCall callType = iota
	contractCall
	erc20Call
)

// ContractData is a helper struct to hold the addresses and ABIs for the
// different contract instances that are subject to testing here.
type CallsData struct {
	sender keyring.Key

	erc20Addr common.Address
	erc20ABI  abi.ABI

	contractAddr common.Address
	contractABI  abi.ABI

	precompileAddr common.Address
	precompileABI  abi.ABI
}

// getCallArgs is a helper function to return the correct call arguments and
// transaction data for a given call type.
func (cd CallsData) getTxAndCallArgs(
	callType callType,
	methodName string,
	args ...interface{},
) (evmtypes.EvmTxArgs, factory.CallArgs) {
	txArgs := evmtypes.EvmTxArgs{}
	callArgs := factory.CallArgs{}

	switch callType {
	case directCall:
		txArgs.To = &cd.precompileAddr
		callArgs.ContractABI = cd.precompileABI
	case contractCall:
		txArgs.To = &cd.contractAddr
		callArgs.ContractABI = cd.contractABI
	case erc20Call:
		txArgs.To = &cd.erc20Addr
		callArgs.ContractABI = cd.erc20ABI
	}

	callArgs.MethodName = methodName
	callArgs.Args = args

	// Setting gas tip cap to zero to have zero gas price.
	txArgs.GasTipCap = new(big.Int).SetInt64(0)

	return txArgs, callArgs
}