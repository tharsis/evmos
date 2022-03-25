package ibc

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
)

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("evmos", "evmospub")
}

func TestGetTransferSenderRecipient(t *testing.T) {
	testCases := []struct {
		name         string
		packet       channeltypes.Packet
		expSender    string
		expRecipient string
		expAmount    sdk.Int
		expError     bool
	}{
		{
			"empty packet",
			channeltypes.Packet{},
			"", "", sdk.ZeroInt(),
			true,
		},
		{
			"invalid packet data",
			channeltypes.Packet{
				Data: ibctesting.MockFailPacketData,
			},
			"", "", sdk.ZeroInt(),
			true,
		},
		{
			"empty FungibleTokenPacketData",
			channeltypes.Packet{
				Data: transfertypes.ModuleCdc.MustMarshalJSON(
					&transfertypes.FungibleTokenPacketData{},
				),
			},
			"", "", sdk.ZeroInt(),
			true,
		},
		{
			"invalid sender",
			channeltypes.Packet{
				Data: transfertypes.ModuleCdc.MustMarshalJSON(
					&transfertypes.FungibleTokenPacketData{
						Sender:   "cosmos1",
						Receiver: "evmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v",
						Amount:   "123456",
					},
				),
			},
			"", "", sdk.ZeroInt(),
			true,
		},
		{
			"invalid recipient",
			channeltypes.Packet{
				Data: transfertypes.ModuleCdc.MustMarshalJSON(
					&transfertypes.FungibleTokenPacketData{
						Sender:   "cosmos1qql8ag4cluz6r4dz28p3w00dnc9w8ueulg2gmc",
						Receiver: "evmos1",
						Amount:   "123456",
					},
				),
			},
			"", "", sdk.ZeroInt(),
			true,
		},
		{
			"valid - cosmos sender, evmos recipient",
			channeltypes.Packet{
				Data: transfertypes.ModuleCdc.MustMarshalJSON(
					&transfertypes.FungibleTokenPacketData{
						Sender:   "cosmos1qql8ag4cluz6r4dz28p3w00dnc9w8ueulg2gmc",
						Receiver: "evmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v",
						Amount:   "123456",
					},
				),
			},
			"evmos1qql8ag4cluz6r4dz28p3w00dnc9w8ueuafmxps",
			"evmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v",
			sdk.NewInt(123456),
			false,
		},
		{
			"valid - evmos sender, cosmos recipient",
			channeltypes.Packet{
				Data: transfertypes.ModuleCdc.MustMarshalJSON(
					&transfertypes.FungibleTokenPacketData{
						Sender:   "evmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v",
						Receiver: "cosmos1qql8ag4cluz6r4dz28p3w00dnc9w8ueulg2gmc",
						Amount:   "123456",
					},
				),
			},
			"evmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v",
			"evmos1qql8ag4cluz6r4dz28p3w00dnc9w8ueuafmxps",
			sdk.NewInt(123456),
			false,
		},
		{
			"valid - osmosis sender, evmos recipient",
			channeltypes.Packet{
				Data: transfertypes.ModuleCdc.MustMarshalJSON(
					&transfertypes.FungibleTokenPacketData{
						Sender:   "osmo1qql8ag4cluz6r4dz28p3w00dnc9w8ueuhnecd2",
						Receiver: "evmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v",
						Amount:   "123456",
					},
				),
			},
			"evmos1qql8ag4cluz6r4dz28p3w00dnc9w8ueuafmxps",
			"evmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v",
			sdk.NewInt(123456),
			false,
		},
	}

	for _, tc := range testCases {
		sender, recipient, _, _, amt, err := GetTransferSenderRecipient(tc.packet)
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
			require.Equal(t, tc.expSender, sender.String())
			require.Equal(t, tc.expRecipient, recipient.String())
			require.Equal(t, tc.expAmount, amt)
		}
	}
}
