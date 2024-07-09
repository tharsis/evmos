// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package app

import (
	"github.com/evmos/evmos/v18/app/eips"
	evmconfig "github.com/evmos/evmos/v18/x/evm/config"
	"github.com/evmos/evmos/v18/x/evm/core/vm"
)

// The init function of the config file allows to setup the global
// configuration for the EVM, modifying the custom ones defined in evmOS.
func init() {
	err := evmconfig.NewEVMConfigurator().
		WithExtendedEips(evmosActivators).
		WithExtendedDefaultExtraEIPs(defaultEnabledEIPs...).
		Configure()
	if err != nil {
		panic(err)
	}
}

var (
	// EvmosActivators defines a map of opcode modifiers associated
	// with a key defining the corresponding EIP.
	evmosActivators = map[int]func(*vm.JumpTable){
		0o000: eips.Enable0000,
		0o001: eips.Enable0001,
		0o002: eips.Enable0002,
	}

	// DefaultEnabledEIPs defines the EIP that should be activated
	// by default and will be merged in the x/evm Params.
	defaultEnabledEIPs = []int64{
		0o000,
		0o001,
		0o002,
	}
)
