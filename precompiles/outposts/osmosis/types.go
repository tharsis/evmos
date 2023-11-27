// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package osmosis

import (
	"encoding/json"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/evmos/v15/utils"
	"golang.org/x/exp/slices"

	cosmosbech32 "github.com/cosmos/cosmos-sdk/types/bech32"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	cmn "github.com/evmos/evmos/v15/precompiles/common"
)

const (
	// MaxSlippagePercentage is the maximum slippage percentage that can be used in the
	// definition of the slippage for the swap.
	MaxSlippagePercentage uint8 = 20
	// MaxWindowSeconds is the maximum number of seconds that can be used in the
	// definition of the slippage for the swap.
	MaxWindowSeconds uint64 = 60
)

const (
	// DefaultOnFailedDelivery is the default value used for the `on_failed_delivery` field.
	DefaultOnFailedDelivery = "do_nothing"
)

const (
	// OsmosisDenom is the base denom in the Osmosis chain.
	OsmosisDenom = "uosmo"
)

// EventSwap is the event type emitted on a Swap transaction.
type EventSwap struct {
	Sender   common.Address
	Input    common.Address
	Output   common.Address
	Amount   *big.Int
	Receiver string
}

// IBCConnction contains information of port and channel of an IBC connection.
type IBCConnection struct {
	PortID    string
	ChannelID string
}

// NewIBCConnection return a new instance of IBCConnection.
func NewIBCConnection(
	portID, ChannelID string,
) IBCConnection {
	return IBCConnection{
		PortID:    portID,
		ChannelID: ChannelID,
	}
}

// TWAP represents a Time-Weighted Average Price configuration.
type TWAP struct {
	// SlippagePercentage specifies the acceptable slippage percentage for a transaction.
	SlippagePercentage uint8 `json:"slippage_percentage"`
	// WindowSeconds defines the duration for which the TWAP is calculated.
	WindowSeconds uint64 `json:"window_seconds"`
}

// Slippage specify how to compute the slippage of the swap. For this version of the outpost
// only the TWAP is allowed.
type Slippage struct {
	TWAP *TWAP `json:"twap"`
}

// OsmosisSwap represents the details for a swap transaction on the Osmosis chain.
//
//nolint:revive
type OsmosisSwap struct {
	// OutputDenom specifies the desired output denomination for the swap.
	OutputDenom string `json:"output_denom"`
	// Slippage represents the TWAP configuration for the swap.
	Slippage *Slippage `json:"slippage"`
	// Receiver is the address of the entity receiving the swapped amount.
	Receiver string `json:"receiver"`
	// OnFailedDelivery specifies the action to be taken in case the swap delivery fails.
	// This can be "do_nothing" or the address on the Osmosis chain that can recover funds
	// in case of errors.
	OnFailedDelivery string `json:"on_failed_delivery"`
	// NextMemo contains any additional memo information for the next operation in a PFM setting.
	NextMemo string `json:"next_memo,omitempty"`
}

// Msg contains the OsmosisSwap details used in the memo relayed to the Osmosis IBC Wasm router.
type Msg struct {
	// OsmosisSwap provides details for a swap transaction.
	OsmosisSwap *OsmosisSwap `json:"osmosis_swap"`
}

// Wasm wraps the message details required for the IBC packet to be valid for the Wasm router.
type Wasm struct {
	// Contract represents the address or identifier of the contract to be called.
	Contract string `json:"contract"`
	// Msg contains the details of the operation to be executed on the contract.
	Msg *Msg `json:"msg"`
}

// Memo wraps the message details for the IBC packet relayed to the Osmosis chain.
type Memo struct {
	Wasm *Wasm `json:"wasm"`
}

// Validate performs basic validation of the IBC memo for the Osmosis outpost.
// This function assumes that memo field is parsed with ParseSwapPacketData, which
// performs data casting ensuring outputDenom cannot be an empty string.
func (m Memo) Validate() error {
	osmosisSwap := m.Wasm.Msg.OsmosisSwap

	if osmosisSwap.OnFailedDelivery == "" {
		return fmt.Errorf(ErrEmptyOnFailedDelivery)
	}

	// Check if account is a valid bech32 address
	if _, err := sdk.AccAddressFromBech32(osmosisSwap.Receiver); err != nil {
		return fmt.Errorf(ErrReceiverAddress, "not a valid evmos address")
	}

	if osmosisSwap.Slippage.TWAP.SlippagePercentage == 0 || osmosisSwap.Slippage.TWAP.SlippagePercentage > MaxSlippagePercentage {
		return fmt.Errorf(ErrSlippagePercentage)
	}

	if osmosisSwap.Slippage.TWAP.WindowSeconds == 0 || osmosisSwap.Slippage.TWAP.WindowSeconds > MaxWindowSeconds {
		return fmt.Errorf(ErrWindowSeconds)
	}

	return nil
}

// RawPacketMetadata is the raw packet metadata used to construct a JSON string.
type RawPacketMetadata struct {
	// The Osmosis outpost IBC memo.
	Memo *Memo `json:"memo"`
}

// CreatePacketWithMemo creates the IBC packet with the memo for the Osmosis
// outpost that can be parsed by the ibc hook middleware, the Wasm hook, on the Osmosis chain.
func CreatePacketWithMemo(
	outputDenom, receiver, contract string,
	slippagePercentage uint8,
	windowSeconds uint64,
	onFailedDelivery, nextMemo string,
) *RawPacketMetadata {
	return &RawPacketMetadata{
		&Memo{
			&Wasm{
				Contract: contract,
				Msg: &Msg{
					OsmosisSwap: &OsmosisSwap{
						OutputDenom: outputDenom,
						Slippage: &Slippage{
							&TWAP{
								SlippagePercentage: slippagePercentage,
								WindowSeconds:      windowSeconds,
							},
						},
						Receiver:         receiver,
						OnFailedDelivery: onFailedDelivery,
						NextMemo:         nextMemo,
					},
				},
			},
		},
	}
}

// ConvertToJSON convert the RawPacketMetadata type into a JSON formatted
// string.
func (r RawPacketMetadata) String() string {
	// Convert the struct to a JSON string
	jsonBytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

// CreateOnFailedDeliveryField is an utility function to create the memo field
// onFailedDelivery. The returned string is the bech32 of the receiver input
// or "do_nothing".
func CreateOnFailedDeliveryField(receiver string) string {
	onFailedDelivery := receiver
	bech32Prefix, addressBytes, err := cosmosbech32.DecodeAndConvert(receiver)
	if err != nil {
		return DefaultOnFailedDelivery
	}
	if bech32Prefix != OsmosisPrefix {
		onFailedDelivery, err = sdk.Bech32ifyAddressBytes(OsmosisPrefix, addressBytes)
		if err != nil {
			return DefaultOnFailedDelivery
		}
	}

	return onFailedDelivery
}

// ValidateInputOutput validate the input and output tokens used in the Osmosis swap.
func ValidateInputOutput(
	inputDenom, outputDenom, stakingDenom string,
	evmosConnection IBCConnection,
) error {
	if outputDenom == inputDenom {
		return fmt.Errorf(ErrInputEqualOutput, inputDenom)
	}

	osmoIBCDenom := utils.ComputeIBCDenom(evmosConnection.PortID, evmosConnection.ChannelID, OsmosisDenom)

	// acceptedTokens are the tokens accepted as input or output of the swap.
	acceptedTokens := []string{stakingDenom, osmoIBCDenom}

	// Check that the input token is aevmos or uosmo.
	if !slices.Contains(acceptedTokens, inputDenom) {
		return fmt.Errorf(ErrTokenNotSupported, acceptedTokens)
	}

	// Check that the output token is aevmos or uosmo.
	if !slices.Contains(acceptedTokens, outputDenom) {
		return fmt.Errorf(ErrTokenNotSupported, acceptedTokens)
	}

	return nil
}

// ConvertToOsmosisRepresentation returns the Osmosis representation of the denom from the Evmos
// representation of aevmos and uosmo. Return an error if the denom is different from one these two.
func ConvertToOsmosisRepresentation(
	denom, stakingDenom string,
	evmosConnection, osmosisConnection IBCConnection,
) (denomOsmosis string, err error) {
	osmoIBCDenom := utils.ComputeIBCDenom(
		evmosConnection.PortID,
		evmosConnection.ChannelID,
		OsmosisDenom,
	)

	switch denom {
	case osmoIBCDenom:
		denomOsmosis = OsmosisDenom
	case stakingDenom:
		denomPrefix := transfertypes.GetPrefixedDenom(
			osmosisConnection.PortID,
			osmosisConnection.ChannelID,
			denom,
		)
		denomTrace := transfertypes.ParseDenomTrace(denomPrefix)
		denomOsmosis = denomTrace.IBCDenom()
	default:
		err = fmt.Errorf(ErrTokenNotSupported, []string{stakingDenom, osmoIBCDenom})
	}
	return denomOsmosis, err
}

// SwapPacketData is an utility structure used to wrap args reiceived by the
// Solidity interface of the Swap function.
type SwapPacketData struct {
	Sender             common.Address
	Input              common.Address
	Output             common.Address
	Amount             *big.Int
	SlippagePercentage uint8
	WindowSeconds      uint64
	SwapReceiver       string
}

// ParseSwapPacketData parses the packet data for the Osmosis swap function.
func ParseSwapPacketData(args []interface{}) (
	swapPacketData SwapPacketData,
	err error,
) {
	if len(args) != 7 {
		return SwapPacketData{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 7, len(args))
	}

	sender, ok := args[0].(common.Address)
	if !ok {
		return SwapPacketData{}, fmt.Errorf(cmn.ErrInvalidType, "sender", common.Address{}, args[0])
	}

	input, ok := args[1].(common.Address)
	if !ok {
		return SwapPacketData{}, fmt.Errorf(cmn.ErrInvalidType, "input", common.Address{}, args[1])
	}

	output, ok := args[2].(common.Address)
	if !ok {
		return SwapPacketData{}, fmt.Errorf(cmn.ErrInvalidType, "output", common.Address{}, args[2])
	}

	amount, ok := args[3].(*big.Int)
	if !ok {
		return SwapPacketData{}, fmt.Errorf(cmn.ErrInvalidType, "amount", big.Int{}, args[3])
	}

	slippagePercentage, ok := args[4].(uint8)
	if !ok {
		return SwapPacketData{}, fmt.Errorf(cmn.ErrInvalidType, "slippagePercentage", uint8(0), args[4])
	}

	windowSeconds, ok := args[5].(uint64)
	if !ok {
		return SwapPacketData{}, fmt.Errorf(cmn.ErrInvalidType, "windowSeconds", uint64(0), args[5])
	}

	receiver, ok := args[6].(string)
	if !ok {
		return SwapPacketData{}, fmt.Errorf(cmn.ErrInvalidType, "receiver", "", args[6])
	}

	swapPacketData = SwapPacketData{
		Sender:             sender,
		Input:              input,
		Output:             output,
		Amount:             amount,
		SlippagePercentage: slippagePercentage,
		WindowSeconds:      windowSeconds,
		SwapReceiver:       receiver,
	}

	return swapPacketData, nil
}
