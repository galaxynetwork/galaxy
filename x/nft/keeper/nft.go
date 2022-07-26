package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

// MintNFT defines a method for minting a new nft
// Note: When the upper module uses this method, it needs to authenticate class
func (k Keeper) MintNFT(ctx sdk.Context, nft types.NFT, recipient sdk.AccAddress) error {
	if exist := k.HasNFT(ctx, nft.BrandId, nft.ClassId, nft.Id); exist {
		return sdkerrors.Wrapf(types.ErrExistNFT, "for brandID: %s, classID: %s, id: %d", nft.BrandId, nft.ClassId, nft.Id)
	}

	if burned := k.BurnedNFT(ctx, nft.BrandId, nft.ClassId, nft.Id); burned {
		return sdkerrors.Wrapf(types.ErrAlreadyBurnedNFT, "for brandID: %s, classID: %s, id: %d", nft.BrandId, nft.ClassId, nft.Id)
	}

	if err := k.setNFT(ctx, nft); err != nil {
		return err
	}

	k.setOwner(ctx, nft.BrandId, nft.ClassId, nft.Id, recipient)

	if err := k.incrSupplyOfClass(ctx, nft.BrandId, nft.ClassId); err != nil {
		return err
	}

	return nil
}

// BurnNFT defines a method for burning a nft from a owner account.
// Note: When the upper module uses this method, it needs to authenticate nft
func (k Keeper) BurnNFT(ctx sdk.Context, brandID, classID string, id uint64) error {
	if exist := k.HasNFT(ctx, brandID, classID, id); !exist {
		return sdkerrors.Wrapf(types.ErrNotFoundNFT, "for brandID: %s, classID: %s, id: %d", brandID, classID, id)
	}

	k.getNFTStore(ctx, brandID, classID).
		Delete(sdk.Uint64ToBigEndian(id))
	k.deleteOwner(ctx, brandID, classID, id)

	if err := k.decrSupplyOfClass(ctx, brandID, classID); err != nil {
		return err
	}

	return nil
}

// UpdateNFT defines a method for updating varUri an exist nft
// Note: When the upper module uses this method, it needs to authenticate nft
func (k Keeper) UpdateNFT(ctx sdk.Context, brandID, classID string, id uint64, varUri string) error {
	nft, exist := k.GetNFT(ctx, brandID, classID, id)
	if !exist {
		return sdkerrors.Wrapf(types.ErrNotFoundNFT, "for brandID: %s, classID: %s, id: %d", brandID, classID, id)
	}

	nft.VarUri = varUri

	if err := k.setNFT(ctx, nft); err != nil {
		return err
	}

	return nil
}

// TransferNFT defines a method for sending a nft to another account.
// Note: When the upper module uses this method, it needs to authenticate nft
func (k Keeper) TransferNFT(ctx sdk.Context, brandID, classID string, id uint64, recipient sdk.AccAddress) error {
	if !k.HasNFT(ctx, brandID, classID, id) {
		return sdkerrors.Wrapf(types.ErrNotFoundNFT, "for brandID: %s, classID: %s, id: %d", brandID, classID, id)
	}

	k.deleteOwner(ctx, brandID, classID, id)
	k.setOwner(ctx, brandID, classID, id, recipient)

	return nil
}

// GenNFT defines a method for generating a new nft with nft sequence number of class.
func (k Keeper) GenNFT(ctx sdk.Context, brandID, classID, uri, varUri string) (nft types.NFT, err error) {
	var id uint64

	id, err = k.getSequenceOfClass(ctx, brandID, classID)
	if err != nil {
		return
	}

	nft = types.NewNFT(id, brandID, classID, uri, varUri)
	return
}

// HasNFT determines whether the specified brandID and classID and id exist
func (k Keeper) HasNFT(ctx sdk.Context, brandID, classID string, id uint64) bool {
	return k.getNFTStore(ctx, brandID, classID).
		Has(sdk.Uint64ToBigEndian(id))
}

// GetNFT returns the nft information of the specified brandID and classID and nftID
func (k Keeper) GetNFT(ctx sdk.Context, brandID, classID string, id uint64) (types.NFT, bool) {
	var nft types.NFT

	bz := k.getNFTStore(ctx, brandID, classID).Get(sdk.Uint64ToBigEndian(id))

	if bz == nil {
		return nft, false
	}

	err := k.cdc.Unmarshal(bz, &nft)

	if err != nil {
		panic(err)
	}

	return nft, true
}

// HasNFT determines whether the specified brandID and classID and id burend
func (k Keeper) BurnedNFT(ctx sdk.Context, brandID, classID string, id uint64) bool {
	if exist := k.HasNFT(ctx, brandID, classID, id); exist {
		return false
	}

	number, err := k.getSequenceOfClass(ctx, brandID, classID)
	if err != nil {
		panic((err))
	}

	// number always waits for the next sequence.
	if id < number {
		return true
	}

	return false
}

func (k Keeper) setNFT(ctx sdk.Context, nft types.NFT) error {
	store := k.getNFTStore(ctx, nft.BrandId, nft.ClassId)

	if bz, err := k.cdc.Marshal(&nft); err != nil {
		return err
	} else {
		store.Set(
			sdk.Uint64ToBigEndian(nft.Id),
			bz,
		)
		return nil
	}
}

func (k Keeper) getNFTStore(ctx sdk.Context, brandID, classID string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.GetNFTStoreKey(brandID, classID))
}
