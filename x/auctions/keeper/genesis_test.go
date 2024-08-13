// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	testkeyring "github.com/evmos/evmos/v19/testutil/integration/evmos/keyring"
	testnetwork "github.com/evmos/evmos/v19/testutil/integration/evmos/network"
	testutiltx "github.com/evmos/evmos/v19/testutil/tx"
	"github.com/evmos/evmos/v19/utils"
	"github.com/evmos/evmos/v19/x/auctions/keeper"
	"github.com/evmos/evmos/v19/x/auctions/types"
)

func TestInitGenesis(t *testing.T) {
	// Define var used in the mutation function
	var existentAccAddress sdk.AccAddress
	moduleAccountBalance := sdk.NewInt(1)

	testCases := []struct {
		name              string
		expPanic          bool
		mutation          func(*types.GenesisState)
		fundModuleAccount bool
		postCheck         func()
	}{
		{
			name:              "valid default",
			expPanic:          false,
			mutation:          func(_ *types.GenesisState) {},
			fundModuleAccount: true,
			postCheck:         func() {},
		},
		{
			name:     "valid with non empty bidder",
			expPanic: false,
			mutation: func(genesis *types.GenesisState) {
				genesis.Bid.Sender = existentAccAddress.String()
				genesis.Bid.Amount.Amount = moduleAccountBalance
			},
			fundModuleAccount: true,
			postCheck:         func() {},
		},
		{
			name:     "invalid non enough balance on auctions module",
			expPanic: true,
			mutation: func(genesis *types.GenesisState) {
				genesis.Bid.Sender = existentAccAddress.String()
				genesis.Bid.Amount.Amount = sdk.NewInt(1)
			},
			fundModuleAccount: false,
			postCheck:         func() {},
		},
		{
			name:     "invalid non empty bidder but zero amount",
			expPanic: true,
			mutation: func(genesis *types.GenesisState) {
				genesis.Bid.Sender = existentAccAddress.String()
			},
			fundModuleAccount: false,
			postCheck:         func() {},
		},
		{
			name:     "invalid sender does not exist",
			expPanic: true,
			mutation: func(genesis *types.GenesisState) {
				genesis.Bid.Sender = sdk.AccAddress(testutiltx.GenerateAddress().Bytes()).String()
			},
			fundModuleAccount: false,
			postCheck:         func() {},
		},
		{
			name:     "invalid empty sender but bid amount not zero",
			expPanic: true,
			mutation: func(genesis *types.GenesisState) {
				genesis.Bid.Amount.Amount = sdk.NewInt(1)
			},
			fundModuleAccount: false,
			postCheck:         func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keyring := testkeyring.New(1)
			network := testnetwork.NewUnitTestNetwork(

				testnetwork.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
			)
			existentAccAddress = keyring.GetKey(0).AccAddr

			genesis := types.DefaultGenesisState()
			tc.mutation(genesis)

			if tc.fundModuleAccount {
				err := network.App.BankKeeper.SendCoinsFromAccountToModule(network.GetContext(), keyring.GetKey(0).AccAddr, types.ModuleName, sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, moduleAccountBalance)))
				require.NoError(t, err, "failed during sending coin to module account")
				require.NoError(t, network.NextBlock())
				auctionModuleAddress := network.App.AccountKeeper.GetModuleAddress(types.ModuleName)
				auctionModuleBalance := network.App.BankKeeper.GetBalance(network.GetContext(), auctionModuleAddress, utils.BaseDenom)
				require.Equal(t, auctionModuleBalance.Amount, moduleAccountBalance)
			}

			if tc.expPanic {
				require.Panics(t, func() {
					keeper.InitGenesis(
						network.GetContext(),
						network.App.AuctionsKeeper,
						*genesis,
					)
				})
			} else {
				require.NotPanics(t, func() {
					keeper.InitGenesis(
						network.GetContext(),
						network.App.AuctionsKeeper,
						*genesis,
					)
				})
			}
		})
	}
}

func TestExportGenesis(t *testing.T) {
	keyring := testkeyring.New(1)
	network := testnetwork.NewUnitTestNetwork(

		testnetwork.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
	)

	exportedGenesis := keeper.ExportGenesis(network.GetContext(), network.App.AuctionsKeeper)
	defaultGenesis := types.DefaultGenesisState()

	require.Equal(t, exportedGenesis.Bid, defaultGenesis.Bid, "expected a different bid")
	require.Equal(t, exportedGenesis.Params, defaultGenesis.Params, "expected different params")
	require.Equal(t, exportedGenesis.Round, defaultGenesis.Round, "expected different round")
}
