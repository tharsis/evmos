// Copyright Tharsis Labs Ltd.(Eidon-chain)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/Eidon-AI/eidon-chain/blob/main/LICENSE)

//go:build test
// +build test

package app

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/Eidon-AI/eidon-chain/v20/app/eips"
	"github.com/Eidon-AI/eidon-chain/v20/utils"
	"github.com/Eidon-AI/eidon-chain/v20/x/evm/core/vm"
	evmtypes "github.com/Eidon-AI/eidon-chain/v20/x/evm/types"
)

// InitializeAppConfiguration allows to setup the global configuration
// for tests within the Eidon-chain EVM. We're not using the sealed flag
// and resetting the configuration to the provided one on every test setup
func InitializeAppConfiguration(chainID string) error {
	coinInfo, found := evmtypes.ChainsCoinInfo[chainID]
	if !found {
		// default to mainnet coin info
		coinInfo = evmtypes.ChainsCoinInfo[utils.MainnetChainID]
	}

	// set the base denom considering if its mainnet or testnet
	if err := setBaseDenom(coinInfo); err != nil {
		return err
	}

	baseDenom, err := sdk.GetBaseDenom()
	if err != nil {
		return err
	}

	ethCfg := evmtypes.DefaultChainConfig(chainID)

	configurator := evmtypes.NewEVMConfigurator()
	// reset configuration to set the new one
	configurator.ResetTestConfig()
	err = configurator.
		WithExtendedEips(eidon-chainActivators).
		WithChainConfig(ethCfg).
		WithEVMCoinInfo(baseDenom, uint8(coinInfo.Decimals)).
		Configure()
	if err != nil {
		return err
	}

	return nil
}

// Eidon-chainActivators defines a map of opcode modifiers associated
// with a key defining the corresponding EIP.
var eidon-chainActivators = map[string]func(*vm.JumpTable){
	"eidon-chain_0": eips.Enable0000,
	"eidon-chain_1": eips.Enable0001,
	"eidon-chain_2": eips.Enable0002,
}

// setBaseDenom registers the display denom and base denom and sets the
// base denom for the chain. The function registers different values based on
// the EvmCoinInfo to allow different configurations in mainnet and testnet.
func setBaseDenom(ci evmtypes.EvmCoinInfo) (err error) {
	// Defer setting the base denom, and capture any potential error from it.
	// So when failing because the denom was already registered, we ignore it and set
	// the corresponding denom to be base denom
	defer func() {
		err = sdk.SetBaseDenom(ci.Denom)
	}()

	if err := sdk.RegisterDenom(ci.DisplayDenom, math.LegacyOneDec()); err != nil {
		return err
	}
	// sdk.RegisterDenom will automatically overwrite the base denom when the new denom units are lower than the current base denom's units.
	return sdk.RegisterDenom(ci.Denom, math.LegacyNewDecWithPrec(1, int64(ci.Decimals)))
}
