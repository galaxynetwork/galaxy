package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

// GetOwner returns the owner information of the specified nft
func (k Keeper) GetOwner(ctx sdk.Context, brandID, classID string, id uint64) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOwnerStoreKey(brandID, classID, id))

	return sdk.AccAddress(bz)
}

func (k Keeper) setOwner(ctx sdk.Context, brandID, classID string, id uint64, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOwnerStoreKey(brandID, classID, id), owner.Bytes())
}

func (k Keeper) deleteOwner(ctx sdk.Context, brandID, classID string, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOwnerStoreKey(brandID, classID, id))
}
