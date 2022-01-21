package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	incentivestypes "github.com/tharsis/evmos/x/incentives/types"
	"github.com/tharsis/evmos/x/inflation/types"
	// poolincentivestypes "github.com/osmosis-labs/osmosis/x/pool-incentives/types"
)

// MintAndAllocateInflation performs inflation minting and allocation
func (k Keeper) MintAndAllocateInflation(ctx sdk.Context, coin sdk.Coin) error {
	// Mint over-allocates by the developer vesting portion, and burn this later
	if err := k.MintCoins(ctx, coin); err != nil {
		panic(err)
	}

	// Allocate minted coins according to allocation proportions
	if err := k.AllocateMintedCoin(ctx, coin); err != nil {
		panic(err)
	}

	if err := k.AllocateTeamVesting(ctx); err != nil {
		panic(err)
	}

	return nil
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoin sdk.Coin) error {
	newCoins := sdk.NewCoins(newCoin)

	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AllocateMintedCoins allocates minted coins from the inflation to external
// modules
func (k Keeper) AllocateMintedCoin(ctx sdk.Context, mintedCoin sdk.Coin) error {
	params := k.GetParams(ctx)
	proportions := params.InflationDistribution

	// Allocate staking rewards into fee collector account
	stakingRewards := sdk.NewCoins(k.GetProportions(ctx, mintedCoin, proportions.StakingRewards))
	err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		k.feeCollectorName,
		stakingRewards,
	)
	if err != nil {
		return err
	}

	// Allocate usage incentives to incentives module account
	usageIncentives := sdk.NewCoins(k.GetProportions(ctx, mintedCoin, proportions.UsageIncentives))
	err = k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		incentivestypes.ModuleName,
		usageIncentives,
	)
	if err != nil {
		return err
	}

	// Allocate community pool inflation to community pool address
	communityPool := sdk.NewCoins(k.GetProportions(ctx, mintedCoin, proportions.CommunityPool))
	err = k.distrKeeper.FundCommunityPool(
		ctx,
		communityPool,
		k.accountKeeper.GetModuleAddress(types.ModuleName),
	)
	if err != nil {
		return err
	}

	// Burn remaining team vesting allocation. These coins are instead allocated
	// from the developer vesting account address, not the inflation module
	// address to comply with taxation policies. We over-minted coins to the
	// inflation module address earlier, in order to allocate according to the
	// allocation proportions.
	moduleAddr := sdk.AccAddress(types.ModuleAddress.Bytes())
	balance := sdk.NewCoins(k.bankKeeper.GetBalance(ctx, moduleAddr, mintedCoin.Denom))
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, balance)

	return err
}

// AllocateTeamVesting allocates the team vesting proportion from the team
// vesting supply
func (k Keeper) AllocateTeamVesting(ctx sdk.Context) error {
	params := k.GetParams(ctx)

	// Check team vesting account balance
	if ok := k.bankKeeper.HasSupply(ctx, params.MintDenom); !ok {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "team vesting account has no supply",
		)
	}

	// Get team vesting provision to allocate from the tharsis account balance. If
	// tharsis account doesn't have sufficient balance to allocate, only allocate
	// remaining coins.
	coin := params.TeamVestingProvision
	tharsisAccount := k.accountKeeper.GetModuleAddress(types.TharsisAccount)
	balance := k.bankKeeper.GetBalance(ctx, tharsisAccount, params.MintDenom)

	if balance.IsLT(coin) {
		coin = balance
	}
	coins := sdk.NewCoins(coin)

	// Allocate teamVesting to community pool when rewards address is empty
	if params.TeamAddress == "" {
		if err := k.distrKeeper.FundCommunityPool(ctx, coins, tharsisAccount); err != nil {
			return err
		}
	}

	// Send coins to teamVestingReceiver account
	teamAddress, err := sdk.AccAddressFromHex(params.TeamAddress)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		tharsisAccount.String(),
		teamAddress,
		coins,
	)

	return err
}

// GetProportions gets the balance of the `MintedDenom` from minted coins and returns coins according to the `AllocationRatio`
func (k Keeper) GetProportions(ctx sdk.Context, mintedCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(mintedCoin.Denom, mintedCoin.Amount.ToDec().Mul(ratio).TruncateInt())
}
