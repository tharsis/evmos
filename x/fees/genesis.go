package fees

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/evmos/v5/x/fees/keeper"
	"github.com/evmos/evmos/v5/x/fees/types"
)

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data types.GenesisState,
) {
	k.SetParams(ctx, data.Params)

	for _, fee := range data.Fees {
		contract := fee.GetContractAddr()
		deployer := fee.GetDeployerAddr()
		withdrawer := fee.GetWithdrawerAddr()

		// Set initial contracts receiving transaction fees
		k.SetFee(ctx, fee)
		k.SetDeployerMap(ctx, deployer, contract)

		if len(withdrawer) != 0 {
			k.SetWithdrawerMap(ctx, withdrawer, contract)
		}
	}
}

// ExportGenesis export module state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
		Fees:   k.GetFees(ctx),
	}
}
