package inflation

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/tharsis/ethermint/types"
	"github.com/tharsis/evmos/x/inflation/keeper"
	"github.com/tharsis/evmos/x/inflation/types"
)

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	ak types.AccountKeeper,
	data types.GenesisState,
) {
	k.SetParams(ctx, data.Params)

	// Ensure inflation module account is set on genesis
	if acc := ak.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		panic("the inflation module account has not been set")
	}

	// Ensure unvested team module account is set on genesis
	if acc := ak.GetModuleAccount(ctx, types.UnvestedTeamAccount); acc == nil {
		panic("the unvested team module account has not been set")
	}

	// ensure team account exists
	if data.Params.TeamAddress != "" {
		teamAddr, err := sdk.AccAddressFromBech32(data.Params.TeamAddress)
		if err != nil {
			panic(fmt.Errorf("invalid team bech32 address: %w", err))
		}

		// TODO is this the right spot to set the account?
		// set inflation team account
		teamAccount := &ethermint.EthAccount{
			BaseAccount: authtypes.NewBaseAccount(teamAddr, nil, 0, 0),
			CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
		}
		ak.SetAccount(ctx, teamAccount)

		if acc := ak.GetAccount(ctx, teamAddr); acc == nil {
			panic(fmt.Errorf("the team account %s has not been set", data.Params.TeamAddress))
		}
	}

	// Set Period
	period := data.Period
	k.SetPeriod(ctx, period)

	// Calculate epoch mint provision
	epochMintProvision := types.CalculateEpochMintProvision(data.Params, period)
	k.SetEpochMintProvision(ctx, epochMintProvision)

	// Mint genesis coins for teamVesting
	amount := sdk.NewInt(200_000_000)
	coins := sdk.NewCoins(sdk.NewCoin(data.Params.MintDenom, amount))
	if err := k.MintGenesisTeamVestingCoins(ctx, coins); err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
		Period: k.GetPeriod(ctx),
	}
}
