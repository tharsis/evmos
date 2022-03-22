package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tharsis/evmos/v2/x/claims/types"
)

var actions = []types.Action{types.ActionVote, types.ActionDelegate, types.ActionEVM, types.ActionIBCTransfer}

// ClaimCoinsForAction removes the claimable amount entry from a claims record
// and transfers it to the user's account
func (k Keeper) ClaimCoinsForAction(
	ctx sdk.Context,
	addr sdk.AccAddress,
	claimsRecord types.ClaimsRecord,
	action types.Action,
	params types.Params,
) (sdk.Int, error) {
	if action == types.ActionUnspecified || action > types.ActionIBCTransfer {
		return sdk.ZeroInt(), sdkerrors.Wrapf(types.ErrInvalidAction, "%d", action)
	}

	// Get claimable amount. Perform a noop if
	// - we are before the start time, after end time, or claims are disabled OR
	// - action already completed and nothing is claimable
	claimableAmount := k.GetClaimableAmountForAction(ctx, claimsRecord, action, params)
	if claimableAmount.IsZero() {
		return sdk.ZeroInt(), nil
	}

	claimedCoins := sdk.Coins{{Denom: params.ClaimsDenom, Amount: claimableAmount}}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, claimedCoins); err != nil {
		return sdk.ZeroInt(), err
	}

	claimsRecord.MarkClaimed(action)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, claimedCoins.String()),
			sdk.NewAttribute(types.AttributeKeyActionType, action.String()),
		),
	})

	if claimsRecord.HasClaimedAll() {
		k.DeleteClaimsRecord(ctx, addr)
	} else {
		k.SetClaimsRecord(ctx, addr, claimsRecord)
	}

	k.Logger(ctx).Info(
		"claimed action",
		"address", addr.String(),
		"action", action.String(),
	)

	return claimableAmount, nil
}

// MergeClaimsRecords merges two independent claims records (sender and
// recipient) into a new instance by summing up the initial claimable amounts
// from both records.

// This method additionally:
//  - Always claims the IBC action, assuming both record haven't claimed it.
//  - Marks an action as claimed for the new instance by performing an XOR operation between the 2 provided records: `merged completed action = sender completed action XOR recipient completed action`
func (k Keeper) MergeClaimsRecords(
	ctx sdk.Context,
	recipient sdk.AccAddress,
	senderClaimsRecord,
	recipientClaimsRecord types.ClaimsRecord,
	params types.Params,
) (mergedRecord types.ClaimsRecord, err error) {
	claimedAmt := sdk.ZeroInt()

	// new total is the sum of the sender and recipient claims records amounts
	totalClaimableAmt := senderClaimsRecord.InitialClaimableAmount.Add(recipientClaimsRecord.InitialClaimableAmount)
	mergedRecord = types.NewClaimsRecord(totalClaimableAmt)

	// iterate over all the available actions and claim the amount if
	// the recipient or sender has completed an action but the other hasn't
	for _, action := range actions {
		senderCompleted := senderClaimsRecord.HasClaimedAction(action)
		recipientCompleted := recipientClaimsRecord.HasClaimedAction(action)

		switch {
		case senderCompleted && recipientCompleted:
			// Both sender and recipient completed the action.
			// Only mark the action as completed
			mergedRecord.MarkClaimed(action)
		case recipientCompleted && !senderCompleted:
			// claim action for sender since the recipient completed it
			amt := k.GetClaimableAmountForAction(ctx, senderClaimsRecord, action, params)
			claimedAmt = claimedAmt.Add(amt)
			mergedRecord.MarkClaimed(action)
		case !recipientCompleted && senderCompleted:
			// claim action for recipient since the sender completed it
			amt := k.GetClaimableAmountForAction(ctx, recipientClaimsRecord, action, params)
			claimedAmt = claimedAmt.Add(amt)
			mergedRecord.MarkClaimed(action)
		case !senderCompleted && !recipientCompleted:
			// Neither sender or recipient completed the action.
			if action != types.ActionIBCTransfer {
				// No-op if the action is not IBC transfer
				continue
			}

			// claim IBC action for both sender and recipient
			amtIBCRecipient := k.GetClaimableAmountForAction(ctx, recipientClaimsRecord, action, params)
			amtIBCSender := k.GetClaimableAmountForAction(ctx, senderClaimsRecord, action, params)
			claimedAmt = claimedAmt.Add(amtIBCRecipient).Add(amtIBCSender)
			mergedRecord.MarkClaimed(action)
		}
	}

	// safety check to prevent error while sending coins from the module escrow balance to the recipient
	if claimedAmt.IsZero() {
		return mergedRecord, nil
	}

	claimedCoins := sdk.Coins{{Denom: params.ClaimsDenom, Amount: claimedAmt}}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, claimedCoins); err != nil {
		return types.ClaimsRecord{}, err
	}

	return mergedRecord, nil
}

// GetClaimableAmountForAction returns claimable amount for a specific action
// done by an address
func (k Keeper) GetClaimableAmountForAction(
	ctx sdk.Context,
	claimsRecord types.ClaimsRecord,
	action types.Action,
	params types.Params,
) sdk.Int {
	// return zero if there are no coins to claim
	if claimsRecord.InitialClaimableAmount.IsNil() || claimsRecord.InitialClaimableAmount.IsZero() {
		return sdk.ZeroInt()
	}

	// check if the entire airdrop has completed. This shouldn't occur since at
	// the end of the airdrop, the EnableClaims param is disabled.
	if !params.IsClaimsActive(ctx.BlockTime()) {
		return sdk.ZeroInt()
	}

	// check if action already completed
	if claimsRecord.HasClaimedAction(action) {
		return sdk.ZeroInt()
	}

	// NOTE: use len(actions)-1 as we don't consider the Unspecified Action
	actionsAmt := int64(len(types.Action_name) - 1)
	initialClaimablePerAction := claimsRecord.InitialClaimableAmount.QuoRaw(actionsAmt)

	// return full claim amount if the elapsed time <= decay start time
	decayStartTime := params.DecayStartTime()
	if !ctx.BlockTime().After(decayStartTime) {
		return initialClaimablePerAction
	}

	// Decrease claimable amount if elapsed time > decay start time.
	// The decrease is calculated proportionally to how much elapsedDeacay period
	// has passed. If you claim early in the decay period, you are entitled to
	// more coins than if you claim at the end of it.
	//
	// Claimable percent = (1 - elapsed decay) x 100
	elapsedDecay := ctx.BlockTime().Sub(decayStartTime)
	elapsedDecayRatio := sdk.NewDec(elapsedDecay.Nanoseconds()).QuoInt64(params.DurationOfDecay.Nanoseconds())
	claimableRatio := sdk.OneDec().Sub(elapsedDecayRatio)

	// calculate the claimable coins, while rounding the decimals
	claimableCoins := initialClaimablePerAction.ToDec().Mul(claimableRatio).RoundInt()
	return claimableCoins
}

// GetUserTotalClaimable returns claimable amount for a specific action done by
// an address at a given block time
func (k Keeper) GetUserTotalClaimable(ctx sdk.Context, addr sdk.AccAddress) sdk.Int {
	totalClaimable := sdk.ZeroInt()

	claimsRecord, found := k.GetClaimsRecord(ctx, addr)
	if !found {
		return sdk.ZeroInt()
	}

	params := k.GetParams(ctx)

	for _, action := range actions {
		claimableForAction := k.GetClaimableAmountForAction(ctx, claimsRecord, action, params)
		totalClaimable = totalClaimable.Add(claimableForAction)
	}

	return totalClaimable
}
