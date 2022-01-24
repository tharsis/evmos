package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/evmos/x/inflation/types"
)

var _ types.QueryServer = Keeper{}

// Params returns params of the mint module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// EpochProvisions returns minter.EpochProvisions of the mint module.
func (k Keeper) EpochProvisions(c context.Context, _ *types.QueryEpochProvisionsRequest) (*types.QueryEpochProvisionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	epochMintProvision, found := k.GetEpochMintProvision(ctx)
	if found {
		panic("the epochMintProvision has was not found")
	}
	return &types.QueryEpochProvisionsResponse{EpochProvisions: epochMintProvision}, nil
}
