// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v18/x/access_control/types"
)

// HasRole returns true if a has a role otherwise false.
func (k Keeper) HasRole(
	ctx sdk.Context,
	contract common.Address,
	role common.Hash,
	account common.Address,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRole)
	return store.Has(types.KeyRole(contract, role, account))
}

// GetRoleAdmin returns the admin role hash that controls `role`.
func (k Keeper) GetRoleAdmin(
	ctx sdk.Context,
	contract common.Address,
	role common.Hash,
) common.Hash {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRoleAdmin)
	roleBz := store.Get(types.KeyRoleAdmin(contract, role))
	if len(roleBz) == 0 {
		return types.RoleDefaultAdmin
	}

	return common.BytesToHash(roleBz)
}

// SetRole sets a role in the store.
func (k Keeper) SetRole(
	ctx sdk.Context,
	contract common.Address,
	role common.Hash,
	account common.Address,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRole)
	store.Set(types.KeyRole(contract, role, account), []byte{0x01})
}

// DeleteRole deletes a role from the store.
func (k Keeper) DeleteRole(
	ctx sdk.Context,
	contract common.Address,
	role common.Hash,
	account common.Address,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRole)
	store.Delete(types.KeyRole(contract, role, account))
}
