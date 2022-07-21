package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

func (k Keeper) SaveClass(ctx sdk.Context, class types.Class) error {
	if k.HasClass(ctx, class.BrandId, class.Id) {
		return sdkerrors.Wrapf(types.ErrExistClass, "for brandID: %s, id: %s", class.BrandId, class.Id)
	}

	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return err
	}

	ctx.KVStore(k.storeKey).
		Set(types.GetClassStoreKey(class.BrandId, class.Id), bz)

	return nil
}

func (k Keeper) SetClass(ctx sdk.Context, class types.Class) error {
	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return err
	}

	ctx.KVStore(k.storeKey).
		Set(types.GetClassStoreKey(class.BrandId, class.Id), bz)

	return nil
}

func (k Keeper) GetClass(ctx sdk.Context, brandID, id string) (types.Class, bool) {
	var class types.Class

	bz := ctx.KVStore(k.storeKey).
		Get(types.GetClassStoreKey(brandID, id))

	if bz == nil {
		return class, false
	}

	err := k.cdc.Unmarshal(bz, &class)
	if err != nil {
		panic(err)
	}

	return class, true
}

func (k Keeper) HasClass(ctx sdk.Context, brandID, id string) bool {
	return ctx.KVStore(k.storeKey).
		Has(types.GetClassStoreKey(brandID, id))
}

func (k Keeper) GetClasses(ctx sdk.Context) (classes types.Classes) {
	iterator := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.KeyPrefixClass,
	).Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var class types.Class
		if err := k.cdc.Unmarshal(iterator.Value(), &class); err != nil {
			panic(err)
		}

		classes = append(classes, class)
	}
	return
}

func (k Keeper) GetClassesOfBrand(ctx sdk.Context, brandID string) (classes types.Classes) {
	iterator := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.GetClassOfBrandPrefix(brandID),
	).Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var class types.Class
		if err := k.cdc.Unmarshal(iterator.Value(), &class); err != nil {
			panic(err)
		}

		classes = append(classes, class)
	}
	return
}

func (k Keeper) getClassPrefixStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixClass)
}

func (k Keeper) getClassOfBrandPrefixStore(ctx sdk.Context, brandID string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.GetClassOfBrandPrefix(brandID))
}
