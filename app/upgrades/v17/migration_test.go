// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package v17_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	v17 "github.com/evmos/evmos/v16/app/upgrades/v17"
	testutils "github.com/evmos/evmos/v16/testutil/integration/evmos/utils"
	erc20types "github.com/evmos/evmos/v16/x/erc20/types"
	"github.com/stretchr/testify/require"
)

func TestConvertToNativeCoinExtensions(t *testing.T) {
	// NOTE: In the setup function we are creating a custom genesis state for the integration network
	// which contains balances for two accounts in different denominations.
	// There is also an ERC-20 smart contract deployed and some tokens minted for each of the accounts.
	// The balances are split between both token representations (IBC coin and ERC-20 token).
	//
	// This genesis state is the starting point to check the migration for the introduction of STR v2.
	// This should ONLY convert native coins for now, which means that the native ERC-20s should be untouched.
	// All native IBC coins should be converted to the native representation and the full balance should be returned
	// both by the bank and the ERC-20 contract.
	// There should be a new ERC-20 EVM extension registered and the ERC-20 contract should be able to be called
	// after being deleted and re-registered as a precompile.
	ts, err := NewConvertERC20CoinsTestSuite()
	require.NoError(t, err, "failed to create test suite")

	res, err := ts.handler.GetTokenPairs()
	require.NoError(t, err, "failed to get token pairs")
	require.Len(t, res.TokenPairs, 1, "unexpected number of token pairs")
	ts.nativeTokenPair = res.TokenPairs[0]

	ts, err = PrepareNetwork(ts)
	require.NoError(t, err, "failed to setup test")

	// NOTE: we are checking the balances of the account before the migration to compare
	// them with the balances after the migration to check that the WEVMOS tokens have been correctly unwrapped.
	balancePreRes, err := ts.handler.GetBalance(ts.keyring.GetAccAddr(testAccount), AEVMOS)
	require.NoError(t, err, "failed to check balances")

	// We check that the minting of tokens for the contract deployer has worked.
	balance, err := GetERC20Balance(ts.factory, ts.keyring.GetPrivKey(erc20Deployer), ts.erc20Contract)
	require.NoError(t, err, "failed to query ERC-20 balance")
	require.Equal(t, mintAmount, balance, "expected different balance after minting ERC-20")

	err = ts.network.NextBlock()
	require.NoError(t, err, "failed to execute block")

	// check that the WEVMOS balance has been increased
	// TODO: port this to the integration test suite
	balance, err = GetERC20Balance(ts.factory, ts.keyring.GetPrivKey(testAccount), ts.wevmosContract)
	require.NoError(t, err, "failed to query ERC-20 balance")
	require.Equal(t, sentWEVMOS.String(), balance.String(), "expected different balance after minting ERC-20")

	logger := ts.network.GetContext().Logger().With("upgrade")

	// Convert the coins back using the upgrade util
	err = v17.ConvertToNativeCoinExtensions(
		ts.network.GetContext(),
		logger,
		ts.network.App.AccountKeeper,
		ts.network.App.BankKeeper,
		ts.network.App.Erc20Keeper,
		ts.wevmosContract,
	)
	require.NoError(t, err, "failed to convert coins")

	err = ts.network.NextBlock()
	require.NoError(t, err, "failed to execute block")

	// We check that the ERC20 converted coins have been added back to the bank balance.
	//
	// NOTE: We are deliberately ONLY checking the balance of the XMPL coin, because the AEVMOS balance was changed
	// through paying transaction fees and they are not affected by the migration.
	err = testutils.CheckBalances(ts.handler, []banktypes.Balance{
		{Address: ts.keyring.GetAccAddr(testAccount).String(), Coins: sdk.NewCoins(sdk.NewInt64Coin(XMPL, 300))},
		{Address: ts.keyring.GetAccAddr(erc20Deployer).String(), Coins: sdk.NewCoins(sdk.NewInt64Coin(XMPL, 200))},
		{Address: bech32WithERC20s.String(), Coins: sdk.NewCoins(sdk.NewInt64Coin(XMPL, 600))},
	})
	require.NoError(t, err, "failed to check balances")

	// We are checking that the WEVMOS tokens have been converted back to the base denomination.
	balancePostRes, err := ts.handler.GetBalance(ts.keyring.GetAccAddr(testAccount), AEVMOS)
	require.NoError(t, err, "failed to check balances")
	// NOTE: can't test for equality because there are transaction fees taken to query the ERC-20 balance ATM.
	require.Greater(t, balancePostRes.Balance.Amount.String(), balancePreRes.Balance.Amount.String(),
		"expected different balance after converting WEVMOS back to unwrapped denom",
	)

	// We check that the token pair was registered as an active precompile.
	evmParams, err := ts.handler.GetEvmParams()
	require.NoError(t, err, "failed to get evm params")
	require.Contains(t, evmParams.Params.ActivePrecompiles, ts.nativeTokenPair.GetERC20Contract().String(),
		"expected token pair precompile to be active",
	)
	require.NotContains(t, evmParams.Params.ActivePrecompiles, ts.nonNativeTokenPair.GetERC20Contract().String(),
		"expected non-native token pair not to be a precompile",
	)

	// NOTE: We check that the ERC20 contract for the native token pair can still be called,
	// even though the original contract code was deleted, and it is now re-deployed
	// as a precompiled contract.
	balance, err = GetERC20BalanceForAddr(
		ts.factory,
		ts.keyring.GetPrivKey(testAccount),
		accountWithERC20s,
		ts.nativeTokenPair.GetERC20Contract(),
	)
	require.NoError(t, err, "failed to query ERC20 balance")
	require.Equal(t, int64(600), balance.Int64(), "expected different balance after converting ERC20")

	// NOTE: We check that the balance of the module address is empty after converting native ERC20s
	balancesRes, err := ts.handler.GetAllBalances(authtypes.NewModuleAddress(erc20types.ModuleName))
	require.NoError(t, err, "failed to get balances")
	require.True(t, balancesRes.Balances.IsZero(), "expected different balance for module account")

	// NOTE: We check that the erc20deployer account still has the minted balance after converting the native ERC20s only.
	balance, err = GetERC20Balance(ts.factory, ts.keyring.GetPrivKey(erc20Deployer), ts.nonNativeTokenPair.GetERC20Contract())
	require.NoError(t, err, "failed to query ERC20 balance")
	require.Equal(t, mintAmount, balance, "expected different balance after converting ERC20")

	// NOTE: We check that there all balance of the WEVMOS contract was withdrawn too.
	balance, err = GetERC20Balance(ts.factory, ts.keyring.GetPrivKey(testAccount), ts.wevmosContract)
	require.NoError(t, err, "failed to query ERC20 balance")
	require.Equal(t, common.Big0.Int64(), balance.Int64(), "expected no WEVMOS left after conversion")
}
