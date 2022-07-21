package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

func (k Keeper) SaveClass(ctx sdk.Context, class types.Class) error {
	if k.HasClass(ctx, class.Id, class.BrandId) {
		return sdkerrors.Wrapf(types.ErrExistClassWithinBrand, "brand_id: %s, class_Id: %s", class.BrandId, class.Id)
	}

	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return err
	}

	classesStore := k.getBrandClassesStore(ctx, class.BrandId)
	classesStore.Set([]byte(class.Id), bz)

	return nil
}

func (k Keeper) SetClass(ctx sdk.Context, class types.Class) error {
	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return err
	}

	classesStore := k.getBrandClassesStore(ctx, class.BrandId)
	classesStore.Set([]byte(class.Id), bz)

	return nil
}

func (k Keeper) GetClass(ctx sdk.Context, id, brandID string) (types.Class, bool) {
	classesStore := k.getBrandClassesStore(ctx, brandID)

	var class types.Class

	bz := classesStore.Get([]byte(id))
	if bz == nil {
		return class, false
	}

	err := k.cdc.Unmarshal(bz, &class)
	if err != nil {
		panic(err)
	}

	return class, true
}

func (k Keeper) HasClass(ctx sdk.Context, id, brandID string) bool {
	classesStore := k.getBrandClassesStore(ctx, brandID)

	return classesStore.Has([]byte(id))
}

func (k Keeper) IteratorClasses(ctx sdk.Context, cb func(class types.Class) (stop bool)) {
	store := k.getClassesStore(ctx)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var class types.Class
		if err := k.cdc.Unmarshal(iterator.Value(), &class); err != nil {
			panic(err)
		}

		if cb(class) {
			break
		}
	}
}

func (k Keeper) IteratorClassesByBrand(ctx sdk.Context, brandID string, cb func(class types.Class) (stop bool)) {
	store := k.getBrandClassesStore(ctx, brandID)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var class types.Class
		if err := k.cdc.Unmarshal(iterator.Value(), &class); err != nil {
			panic(err)
		}

		if cb(class) {
			break
		}
	}
}

func (k Keeper) GetClasses(ctx sdk.Context) (classes types.Classes) {
	k.IteratorClasses(ctx, func(class types.Class) (stop bool) {
		classes = append(classes, class)
		return false
	})

	return
}

func (k Keeper) GetClassesByBrand(ctx sdk.Context, brandID string) (classes types.Classes) {
	k.IteratorClassesByBrand(ctx, brandID, func(class types.Class) (stop bool) {
		classes = append(classes, class)
		return false
	})

	return
}

func (k Keeper) getClassesStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.PrefixClassKey)
}

func (k Keeper) getBrandClassesStore(ctx sdk.Context, brandID string) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetClassStoreKey(brandID))
}
