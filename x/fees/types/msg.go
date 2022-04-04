package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/tharsis/ethermint/types"
)

var (
	_ sdk.Msg = &MsgRegisterDevFeeInfo{}
	_ sdk.Msg = &MsgCancelDevFeeInfo{}
	_ sdk.Msg = &MsgUpdateDevFeeInfo{}
)

const (
	TypeMsgRegisterDevFeeInfo = "register_fee_contract"
	TypeMsgCancelDevFeeInfo   = "cancel_fee_contract"
	TypeMsgUpdateDevFeeInfo   = "update_fee_contract"
)

// NewMsgRegisterDevFeeInfo creates new instance of MsgRegisterDevFeeInfo
func NewMsgRegisterDevFeeInfo(
	contract common.Address,
	deployer sdk.AccAddress,
	withdraw sdk.AccAddress,
	nonces []uint64,
) *MsgRegisterDevFeeInfo {
	return &MsgRegisterDevFeeInfo{
		ContractAddress: contract.String(),
		DeployerAddress: deployer.String(),
		WithdrawAddress: withdraw.String(),
		Nonces:          nonces,
	}
}

// Route returns the name of the module
func (msg MsgRegisterDevFeeInfo) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterDevFeeInfo) Type() string { return TypeMsgRegisterDevFeeInfo }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterDevFeeInfo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if err := ethermint.ValidateAddress(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid contract address %s", msg.ContractAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.WithdrawAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid withdraw address address %s", msg.WithdrawAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterDevFeeInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterDevFeeInfo) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.DeployerAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{from}
}

// NewMsgClawbackcreates new instance of MsgClawback. The dest_address may be
// nil - defaulting to the funder.
func NewMsgCancelDevFeeInfo(deployer sdk.AccAddress, contract string) *MsgCancelDevFeeInfo {
	return &MsgCancelDevFeeInfo{
		ContractAddress: contract,
		DeployerAddress: deployer.String(),
	}
}

// Route returns the message route for a MsgCancelDevFeeInfo.
func (msg MsgCancelDevFeeInfo) Route() string { return RouterKey }

// Type returns the message type for a MsgCancelDevFeeInfo.
func (msg MsgCancelDevFeeInfo) Type() string { return TypeMsgCancelDevFeeInfo }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelDevFeeInfo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if err := ethermint.ValidateAddress(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid contract address %s", msg.ContractAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelDevFeeInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCancelDevFeeInfo) GetSigners() []sdk.AccAddress {
	funder, err := sdk.AccAddressFromBech32(msg.DeployerAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{funder}
}

// NewMsgUpdateDevFeeInfo creates new instance of MsgUpdateDevFeeInfo
func NewMsgUpdateDevFeeInfo(
	deployer sdk.AccAddress,
	contract string,
	withdraw sdk.AccAddress,
) *MsgUpdateDevFeeInfo {
	return &MsgUpdateDevFeeInfo{
		DeployerAddress: deployer.String(),
		ContractAddress: contract,
		WithdrawAddress: withdraw.String(),
	}
}

// Route returns the name of the module
func (msg MsgUpdateDevFeeInfo) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateDevFeeInfo) Type() string { return TypeMsgUpdateDevFeeInfo }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateDevFeeInfo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if err := ethermint.ValidateAddress(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid contract address %s", msg.ContractAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.WithdrawAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid withdraw address address %s", msg.WithdrawAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgUpdateDevFeeInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateDevFeeInfo) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.DeployerAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{from}
}
