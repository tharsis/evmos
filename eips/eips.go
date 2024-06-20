// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package eips

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/evmos/evmos/v18/x/evm/core/vm"
)

// EvmosActivators defines a map of opcode modifiers associated
// with a key defining the corresponding EIP.
var EvmosActivators = map[int]func(*vm.JumpTable){
	0o00000: enable0000,
	0o00001: enable0001,
	0o00002: enable0002,
}

// DefaultEnabledEvmosEIPs defines the EIP that should be activated
// by default and will be merged in the x/evm Params.
var DefaultEnabledEvmosEIPs = []int64{
	0o00000,
	0o00001,
	0o00002,
}

// enable0000 contains the logic to modify the CREATE and CREATE2 opcodes
// constant gas value.
// TODO: define the multiplier.
func enable0000(jt *vm.JumpTable) {
	// Gas cost changes
	jt[vm.CREATE].SetConstantGas(params.CreateGas * 2)
	jt[vm.CREATE2].SetConstantGas(params.CreateGas * 2)
}

// enable0001 contains the logic to modify the CALL opcode
// constant gas value.
// TODO: define the multiplier.
func enable0001(jt *vm.JumpTable) {
	// Gas cost changes
	jt[vm.CALL].SetConstantGas(params.CallGasEIP150 * 2)
}

// enable0002 contains the logic to modify the SSTORE opcode
// constant gas value.
// TODO: define the multiplier.
func enable0002(jt *vm.JumpTable) {
	// Gas cost changes
	jt[vm.SSTORE].SetConstantGas(params.SstoreSetGas * 2)
}
