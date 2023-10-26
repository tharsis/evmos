// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
//
// Osmosis package contains the logic of the Osmosis outpost on the Evmos chain.
// This outpost uses the ics20 precompile to relay IBC packets to the Osmosis
// chain, targeting the XCSV
package osmosis

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v15/precompiles/ics20"

	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// SwapMethod is the name of the swap method
	SwapMethod = "swap"
	// SwapAction is the action name needed in the memo field
	SwapAction = "Swap"
)

const (
	// NextMemo is the memo to use after the swap of the token in the IBC packet
	// built on the Osmosis chain. In the alpha version of the outpost this is
	// an empty string that will not be included in the XCS contract payload.
	NextMemo = ""

	// TODO: XCSContract is the swap contract on the Osmosis chain
	XCSContract = "placeholder"

	// OsmosisPrefix is the prefix for osmosis addresses
	OsmosisPrefix = "osmo"
)

// Swap is a transaction that swap tokens on the Osmosis chain using
// an ICS20 transfer with a custom memo field to trigger the XCS V2 contract.
func (p Precompile) Swap(
	ctx sdk.Context,
	origin common.Address,
	stateDB vm.StateDB,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	sender, input, output, amount, slippagePercentage, windowSeconds, receiver, err := ParseSwapPacketData(args)
	if err != nil {
		return nil, err
	}

	// The provided sender address should always be equal to the origin address.
	// In case the contract caller address is the same as the sender address provided,
	// update the sender address to be equal to the origin address.
	// Otherwise, if the provided sender address is different from the origin address,
	// return an error because is a forbidden operation
	sender, err = ics20.CheckOriginAndSender(contract, origin, sender)
	if err != nil {
		return nil, err
	}

	inputTokenPairID := p.erc20Keeper.GetERC20Map(ctx, input)
	inputTokenPair, found := p.erc20Keeper.GetTokenPair(ctx, inputTokenPairID)
	if !found {
		return nil, fmt.Errorf(ErrTokenPairNotFound, input)
	}
	inputDenom := inputTokenPair.Denom

	outputTokenPairID := p.erc20Keeper.GetERC20Map(ctx, output)
	outputTokenPair, found := p.erc20Keeper.GetTokenPair(ctx, outputTokenPairID)
	if !found {
		return nil, fmt.Errorf(ErrTokenPairNotFound, output)
	}
	outputDenom := outputTokenPair.Denom

	// We need the bonded denom just for the outpost alpha version where the
	// the only two inputs allowed are aevmos and uosmo.
	bondDenom := p.stakingKeeper.GetParams(ctx).BondDenom

	// If the receiver has not the prefix "osmo", we should compute its address
	// in the Osmosis chain as a recovery address for the contract. This address
	// is computed on the outpost for the alpha version just to be sure that it
	// is provided in the payload.
	onFailedDelivery := receiver
	bech32Prefix, address, err := bech32.DecodeAndConvert(receiver)
	if err != nil {
		return nil, err
	}
	if bech32Prefix != OsmosisPrefix {
		onFailedDelivery, err = sdk.Bech32ifyAddressBytes(OsmosisDenom, address)
		if err != nil {
			return nil, err
		}
	}

	packet := CreatePacketWithMemo(
		outputDenom, receiver, XCSContract, slippagePercentage, windowSeconds, onFailedDelivery, NextMemo,
	)
	err = packet.Memo.Validate()
	if err != nil {
		return nil, err
	}
	packetString := packet.String()

	coin := sdk.Coin{Denom: inputDenom, Amount: sdk.NewIntFromBigInt(amount)}
	msg, err := ics20.CreateAndValidateMsgTransfer(
		p.portID,
		p.channelID,
		coin,
		sdk.AccAddress(sender.Bytes()).String(),
		receiver,
		p.timeoutHeight,
		p.timeoutTimestamp,
		packetString,
	)

	return nil, nil
}

// GetTokenDenom returns the denom associated to the tokenAddress from the
// erc20 store. Returns an error if the TokenPair associated to the tokenAddress
// is not found.
func (p Precompile) GetTokenDenom(ctx sdk.Context, tokenAddress common.Address) (string, error) {
	TokenPairID := p.erc20Keeper.GetERC20Map(ctx, tokenAddress)
	TokenPair, found := p.erc20Keeper.GetTokenPair(ctx, TokenPairID)
	if !found {
		return "", fmt.Errorf(ErrTokenPairNotFound, tokenAddress)
	}

	return TokenPair.Denom, nil
}
