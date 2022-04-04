package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/tharsis/evmos/v3/x/fees/types"
)

var _ types.MsgServer = &Keeper{}

// RegisterDevFeeInfo registers a contract to receive transaction fees
func (k Keeper) RegisterDevFeeInfo(
	goCtx context.Context,
	msg *types.MsgRegisterDevFeeInfo,
) (*types.MsgRegisterDevFeeInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.isEnabled(ctx) {
		return nil, sdkerrors.Wrapf(types.ErrInternalFee, "fees module is not enabled")
	}

	contract := common.HexToAddress(msg.ContractAddress)
	if k.IsFeeRegistered(ctx, contract) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "contract %s is already registered", contract)
	}

	var withdrawal *sdk.AccAddress
	deployer, _ := sdk.AccAddressFromBech32(msg.DeployerAddress)
	derivedContractAddr := common.BytesToAddress(deployer)

	if msg.WithdrawAddress != "" {
		_withdrawal, _ := sdk.AccAddressFromBech32(msg.WithdrawAddress)
		withdrawal = &_withdrawal
	}

	for _, nonce := range msg.Nonces {
		derivedContractAddr = crypto.CreateAddress(derivedContractAddr, nonce)
	}

	if contract != derivedContractAddr {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrorInvalidSigner,
			"%s not contract deployer or wrong nonce", msg.DeployerAddress,
		)
	}

	// contract must already be deployed, to avoid spam registrations
	contractAccount := k.evmKeeper.GetAccountWithoutBalance(ctx, contract)
	if contractAccount == nil || !contractAccount.IsContract() {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "contract %s has no code", contract)
	}

	k.SetFee(ctx, contract, deployer, withdrawal)
	k.SetFeeInverse(ctx, deployer, contract)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeRegisterDevFeeInfo,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
				sdk.NewAttribute(types.AttributeKeyWithdrawAddress, msg.WithdrawAddress),
			),
		},
	)

	return &types.MsgRegisterDevFeeInfoResponse{}, nil
}

// CancelDevFeeInfo deletes the fee for a contract
func (k Keeper) CancelDevFeeInfo(
	goCtx context.Context,
	msg *types.MsgCancelDevFeeInfo,
) (*types.MsgCancelDevFeeInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.isEnabled(ctx) {
		return nil, sdkerrors.Wrapf(types.ErrInternalFee, "fees module is not enabled")
	}

	deployerAddress, found := k.GetDeployer(ctx, common.HexToAddress(msg.ContractAddress))
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "contract %s is not registered", msg.ContractAddress)
	}

	if msg.DeployerAddress != deployerAddress.String() {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not the contract deployer", msg.DeployerAddress)
	}

	k.DeleteFee(ctx, common.HexToAddress(msg.ContractAddress))
	k.DeleteFeeInverse(ctx, deployerAddress)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeCancelDevFeeInfo,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
			),
		},
	)

	return &types.MsgCancelDevFeeInfoResponse{}, nil
}

// UpdateDevFeeInfo updates the withdraw address for a contract
func (k Keeper) UpdateDevFeeInfo(
	goCtx context.Context,
	msg *types.MsgUpdateDevFeeInfo,
) (*types.MsgUpdateDevFeeInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.isEnabled(ctx) {
		return nil, sdkerrors.Wrapf(types.ErrInternalFee, "fees module is not enabled")
	}

	contractAddress := common.HexToAddress(msg.ContractAddress)
	deployerAddress, found := k.GetDeployer(ctx, contractAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "contract %s is not registered", msg.ContractAddress)
	}

	if msg.DeployerAddress != deployerAddress.String() {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not the contract deployer", msg.DeployerAddress)
	}

	withdrawalAddress, _ := sdk.AccAddressFromBech32(msg.WithdrawAddress)
	k.SetWithdrawal(ctx, contractAddress, withdrawalAddress)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateDevFeeInfo,
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyWithdrawAddress, msg.WithdrawAddress),
			),
		},
	)

	return &types.MsgUpdateDevFeeInfoResponse{}, nil
}
