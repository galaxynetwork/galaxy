package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/galaxies-labs/galaxy/x/mint/types"
)

func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MinterKey)

	if bz == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(bz, &minter)

	return minter
}

func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}
