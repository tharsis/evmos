package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// Hooks wrapper struct for fees keeper
type Hooks struct {
	k Keeper
}

var _ evmtypes.EvmHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing. After each successful
// interaction with an incentivized contract, the owner's GasUsed is
// added to its gasMeter.
func (h Hooks) PostTxProcessing(ctx sdk.Context, from common.Address, contract *common.Address, receipt *ethtypes.Receipt) error {
	// check if the fees are globally enabled
	params := h.k.GetParams(ctx)
	if !params.EnableFees {
		return nil
	}

	// If theres no fees registered for the contract, do nothing
	if contract == nil || !h.k.IsFeeRegistered(ctx, *contract) {
		return nil
	}
	cfg, err := h.k.evmKeeper.EVMConfig(ctx)
	// TODO do something with the error?
	if err != nil {
		return nil
	}

	feeContract, ok := h.k.GetFee(ctx, *contract)
	if !ok {
		return nil
	}
	waddr := common.HexToAddress(feeContract.WithdrawAddress)
	withdrawAddr := sdk.AccAddress(waddr.Bytes())
	distrFees := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), cfg.BaseFee)
	developerFee := new(big.Int).Mul(distrFees, new(big.Int).SetUint64(h.k.GetParams(ctx).DeveloperPercentage))
	developerFee = new(big.Int).Quo(developerFee, big.NewInt(100))

	return h.addFeesToOwner(ctx, *contract, withdrawAddr, developerFee, cfg.Params.EvmDenom)
}

// addGasToParticipant adds gasUsed to a participant's gas meter's cumulative
// gas used
func (h Hooks) addFeesToOwner(
	ctx sdk.Context,
	contract common.Address,
	withdrawAddr sdk.AccAddress,
	fees *big.Int,
	denom string,
) error {
	coins := sdk.Coins{sdk.NewCoin(denom, sdk.NewIntFromBigInt(fees))}
	err := h.k.bankKeeper.SendCoinsFromModuleToAccount(ctx, h.k.feeCollectorName, withdrawAddr, coins)
	if err != nil {
		err = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "fee collector account failed to distribute developer fees: %s", err.Error())
		return sdkerrors.Wrapf(err, "failed to distribute %s fees", fees.String())
	}
	return nil
}
