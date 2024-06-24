package keeper_test

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/evmos/v18/utils"
	"github.com/evmos/evmos/v18/x/evm/keeper"
	"github.com/evmos/evmos/v18/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v18/x/evm/types"

	"github.com/ethereum/go-ethereum/common"
)

func (suite *KeeperTestSuite) TestWithChainID() {
	testCases := []struct {
		name       string
		chainID    string
		expChainID int64
		expPanic   bool
	}{
		{
			"fail - chainID is empty",
			"",
			0,
			true,
		},
		{
			"success - Evmos mainnet chain ID",
			"evmos_9001-2",
			9001,
			false,
		},
		{
			"success - Evmos testnet chain ID",
			"evmos_9000-4",
			9000,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			keeper := keeper.Keeper{}
			ctx := suite.network.GetContext().WithChainID(tc.chainID)

			if tc.expPanic {
				suite.Require().Panics(func() {
					keeper.WithChainID(ctx)
				})
			} else {
				suite.Require().NotPanics(func() {
					keeper.WithChainID(ctx)
					suite.Require().Equal(tc.expChainID, keeper.ChainID().Int64())
				})
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBaseFee() {
	testCases := []struct {
		name            string
		enableLondonHF  bool
		enableFeemarket bool
		expectBaseFee   *big.Int
	}{
		{"not enable london HF, not enable feemarket", false, false, nil},
		{"enable london HF, not enable feemarket", true, false, big.NewInt(0)},
		{"enable london HF, enable feemarket", true, true, big.NewInt(1000000000)},
		{"not enable london HF, enable feemarket", false, true, nil},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.enableFeemarket = tc.enableFeemarket
			suite.enableLondonHF = tc.enableLondonHF
			suite.SetupTest()
			suite.Require().NoError(suite.network.App.EvmKeeper.BeginBlock(suite.network.GetContext()))
			params := suite.network.App.EvmKeeper.GetParams(suite.network.GetContext())
			ethCfg := params.ChainConfig.EthereumConfig(suite.network.App.EvmKeeper.ChainID())
			baseFee := suite.network.App.EvmKeeper.GetBaseFee(suite.network.GetContext(), ethCfg)
			suite.Require().Equal(tc.expectBaseFee, baseFee)
		})
	}
	suite.enableFeemarket = false
	suite.enableLondonHF = true
}

func (suite *KeeperTestSuite) TestGetAccountStorage() {
	var ctx sdk.Context
	testCases := []struct {
		name       string
		malleate   func() map[uint64]int
		expStorage bool
	}{
		{
			"Only one account that's not a contract (no storage)",
			func() map[uint64]int {
				expRes := make(map[uint64]int)
				i := 0
				// NOTE: here we're removing all accounts except for one
				suite.network.App.AccountKeeper.IterateAccounts(ctx, func(account sdk.AccountI) bool {
					defer func() { i++ }()
					if i == 0 {
						expRes[account.GetAccountNumber()] = 0
						return false
					}
					suite.network.App.AccountKeeper.RemoveAccount(ctx, account)
					return false
				})
				return expRes
			},
			false,
		},
		{
			"Two accounts - one contract (with storage), one wallet",
			func() map[uint64]int {
				expRes := make(map[uint64]int)
				supply := big.NewInt(100)
				suite.DeployTestContract(suite.T(), ctx, suite.keyring.GetAddr(0), supply)
				i := 0
				suite.network.App.AccountKeeper.IterateAccounts(ctx, func(account sdk.AccountI) bool {
					defer func() { i++ }()

					var storage evmtypes.Storage
					hexAddr := utils.CosmosToEthAddr(account.GetAddress())
					if suite.network.App.EvmKeeper.IsContract(ctx, hexAddr) {
						storage = suite.network.App.EvmKeeper.GetAccountStorage(ctx, hexAddr)
					}
					expRes[account.GetAccountNumber()] = len(storage)

					if i == 0 || len(storage) > 0 {
						return false
					}

					suite.network.App.AccountKeeper.RemoveAccount(ctx, account)
					return false
				})
				return expRes
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var passed bool
			suite.SetupTest()
			ctx = suite.network.GetContext()
			expRes := tc.malleate()

			suite.network.App.AccountKeeper.IterateAccounts(ctx, func(account sdk.AccountI) bool {
				addr := utils.CosmosToEthAddr(account.GetAddress())
				storage := suite.network.App.EvmKeeper.GetAccountStorage(ctx, addr)

				storageEntriesCount := len(storage)
				expCount, ok := expRes[account.GetAccountNumber()]
				suite.Require().True(ok)
				suite.Require().Equal(expCount, storageEntriesCount)
				if !tc.expStorage {
					if storageEntriesCount > 0 {
						println("Expected no storage entries, but got some")
						passed = false
						return true
					}
					passed = true
				}
				if tc.expStorage && storageEntriesCount > 0 {
					passed = true
				}
				return false
			})
			suite.Require().True(passed)
		})
	}
}

func (suite *KeeperTestSuite) TestGetAccountOrEmpty() {
	ctx := suite.network.GetContext()
	empty := statedb.Account{
		Balance:  new(big.Int),
		CodeHash: evmtypes.EmptyCodeHash,
	}

	supply := big.NewInt(100)
	contractAddr := suite.DeployTestContract(suite.T(), ctx, suite.keyring.GetAddr(0), supply)

	testCases := []struct {
		name     string
		addr     common.Address
		expEmpty bool
	}{
		{
			"unexisting account - get empty",
			common.Address{},
			true,
		},
		{
			"existing contract account",
			contractAddr,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			res := suite.network.App.EvmKeeper.GetAccountOrEmpty(ctx, tc.addr)
			if tc.expEmpty {
				suite.Require().Equal(empty, res)
			} else {
				suite.Require().NotEqual(empty, res)
			}
		})
	}
}
