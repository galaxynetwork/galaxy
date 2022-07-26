package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/internal/conv"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

// GetOwner returns the owner information of the specified nft
func (k Keeper) GetOwner(ctx sdk.Context, brandID, classID string, id uint64) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOwnerStoreKey(brandID, classID, id))

	return sdk.AccAddress(bz)
}

// GetOwner returns the all nfts by owner
func (k Keeper) GetNFTsByOwner(ctx sdk.Context, owner sdk.AccAddress) (nfts types.NFTs) {
	store := k.getNFTStoreByOwner(ctx, owner)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		brandID, classId, id, err := types.ParseNFTUniqueID(conv.UnsafeBytesToStr(iterator.Value()))
		if err != nil {
			panic("invalid nftUniqueID stored in nftOfClassByOwner")
		}

		nft, exist := k.GetNFT(ctx, brandID, classId, id)
		if !exist {
			panic("unexpected nft is stored in nftOfClassByOwner store")
		}
		nfts = append(nfts, nft)
	}

	return
}

func (k Keeper) setOwner(ctx sdk.Context, brandID, classID string, id uint64, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOwnerStoreKey(brandID, classID, id), owner.Bytes())

	classStore := k.getClassStoreByOwner(ctx, owner, brandID, classID)
	classStore.Set(sdk.Uint64ToBigEndian(id), []byte(types.GetNFTUniqueID(brandID, classID, id)))
}

func (k Keeper) deleteOwner(ctx sdk.Context, brandID, classID string, id uint64, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOwnerStoreKey(brandID, classID, id))

	classStore := k.getClassStoreByOwner(ctx, owner, brandID, classID)
	classStore.Delete(sdk.Uint64ToBigEndian(id))
}

func (k Keeper) getClassStoreByOwner(ctx sdk.Context, owner sdk.AccAddress, brandID, classID string) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetNFTOfClassByOwnerStoreKey(owner, brandID, classID))
}

func (k Keeper) getNFTStoreByOwner(ctx sdk.Context, owner sdk.AccAddress) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetPrefixNFTOfClassByOwnerKey(owner))
}
