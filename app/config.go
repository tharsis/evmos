// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

//go:build !test
// +build !test

package app

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/evmos/v20/app/eips"
	"github.com/evmos/evmos/v20/types"
	"github.com/evmos/evmos/v20/utils"
	"github.com/evmos/evmos/v20/x/evm/core/vm"
	evmtypes "github.com/evmos/evmos/v20/x/evm/types"
)

var (
	sealed = false
	// firstCall is used to discriminate between the call to the InitializeAppConfiguration for the CLI client or for the node operator.
	firstCall = true
)

// InitializeAppConfiguration allows to setup the global configuration
// for the Evmos EVM.
func InitializeAppConfiguration(chainID string) error {
	// When calling any CLI command, it creates a tempApp inside RootCmdHandler that will be overwritten later if needed.
	// The configurator can be set with a dirty state only once
	if chainID == "" {
		if firstCall {
			firstCall = false
			return nil
		} else {
			panic("calling configurator twice with invalid chainID")
		}
	}

	if sealed {
		return nil
	}

	// set the base denom considering if its mainnet or testnet
	if err := setBaseDenomWithChainID(chainID); err != nil {
		return err
	}

	baseDenom, err := sdk.GetBaseDenom()
	if err != nil {
		return err
	}

	ethCfg := evmtypes.DefaultChainConfig(chainID)

	err = evmtypes.NewEVMConfigurator().
		WithExtendedEips(evmosActivators).
		WithChainConfig(ethCfg).
		WithEVMCoinInfo(baseDenom, evmtypes.EighteenDecimals).
		Configure()
	if err != nil {
		return err
	}

	sealed = true

	// if the first call was made with the correct chainID (a call without Cobra CLI as entrypoint), it no longer accept empty string as a valid chainID
	firstCall = false

	return nil
}

// EvmosActivators defines a map of opcode modifiers associated
// with a key defining the corresponding EIP.
var evmosActivators = map[string]func(*vm.JumpTable){
	"evmos_0": eips.Enable0000,
	"evmos_1": eips.Enable0001,
	"evmos_2": eips.Enable0002,
}

// setBaseDenomWithChainID registers the display denom and base denom and sets the
// base denom for the chain. The function registers different values based on
// the chainID to allow different configurations in mainnet and testnet.
func setBaseDenomWithChainID(chainID string) error {
	if utils.IsTestnet(chainID) {
		return setTestnetBaseDenom()
	}
	return setMainnetBaseDenom()
}

func setTestnetBaseDenom() error {
	if err := sdk.RegisterDenom(types.DisplayDenomTestnet, math.LegacyOneDec()); err != nil {
		return err
	}
	if err := sdk.RegisterDenom(types.BaseDenomTestnet, math.LegacyNewDecWithPrec(1, types.BaseDenomUnit)); err != nil {
		return err
	}
	return sdk.SetBaseDenom(types.BaseDenomTestnet)
}

func setMainnetBaseDenom() error {
	if err := sdk.RegisterDenom(types.DisplayDenom, math.LegacyOneDec()); err != nil {
		return err
	}
	if err := sdk.RegisterDenom(types.BaseDenom, math.LegacyNewDecWithPrec(1, types.BaseDenomUnit)); err != nil {
		return err
	}
	return sdk.SetBaseDenom(types.BaseDenom)
}
