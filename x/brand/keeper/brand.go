package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

// HasBrand determines whether the specified classID exist
func (keeper Keeper) HasBrand(ctx sdk.Context, brandID string) bool {
	prefix := prefix.NewStore(ctx.KVStore(keeper.storeKey), types.KeyPrefixBrand)

	return prefix.Has(types.GetBrandKey(brandID))
}

// SetBrand defines a method for set an brand in the store.
func (keeper Keeper) SetBrand(ctx sdk.Context, brand types.Brand) error {
	prefix := prefix.NewStore(ctx.KVStore(keeper.storeKey), types.KeyPrefixBrand)

	bz, err := keeper.MarshalBrand(brand)
	if err != nil {
		return err
	}

	prefix.Set(types.GetBrandKey(brand.Id), bz)

	return nil
}

// GetBrand defines a method for returning a existing brand
func (keeper Keeper) GetBrand(ctx sdk.Context, brandID string) (types.Brand, bool) {
	prefix := prefix.NewStore(ctx.KVStore(keeper.storeKey), types.KeyPrefixBrand)

	var brand types.Brand

	bz := prefix.Get(types.GetBrandKey(brandID))

	if bz == nil {
		return brand, false
	}

	if err := keeper.UnmarshalBrand(bz, &brand); err != nil {
		panic(fmt.Errorf("stored brand unmarshalling error: %v", err))
	}

	return brand, true
}

func (keeper Keeper) IterateBrands(ctx sdk.Context, cb func(brand types.Brand) (stop bool)) {

	prefix := prefix.NewStore(ctx.KVStore(keeper.storeKey), types.KeyPrefixBrand)

	iterator := prefix.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var brand types.Brand
		if err := keeper.UnmarshalBrand(iterator.Value(), &brand); err != nil {
			panic(fmt.Errorf("stored brand unmarshalling error: %v", err))
		}

		if cb(brand) {
			break
		}
	}
}
func (keeper Keeper) IterateBrandsByOwner(ctx sdk.Context, owner string, cb func(brand types.Brand) (stop bool)) {
	acc, _ := sdk.AccAddressFromBech32(owner)
	ownerStore := keeper.getBrandByOwnerStore(ctx, acc)
	iterator := ownerStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if brand, exist := keeper.GetBrand(ctx, string(iterator.Value()[:])); !exist {
			panic(fmt.Errorf("unexpected brand is stored in brand_by_owner prefix"))
		} else {
			if cb(brand) {
				break
			}
		}

	}
}

// GetBrands defines a method for returning all brands
func (keeper Keeper) GetBrands(ctx sdk.Context) (brands types.Brands) {
	keeper.IterateBrands(ctx, func(brand types.Brand) (stop bool) {
		brands = append(brands, brand)
		return false
	})
	return brands
}

// GetBrandsByOwner defines a method for returning all brands by owner
func (keeper Keeper) GetBrandsByOwner(ctx sdk.Context, owner string) (brands types.Brands) {
	keeper.IterateBrandsByOwner(ctx, owner, func(brand types.Brand) (stop bool) {
		brands = append(brands, brand)
		return false
	})
	return brands
}

// MarshalBrand defines a method for protobuf serializes brand
func (keeper Keeper) MarshalBrand(brand types.Brand) ([]byte, error) {
	if bz, err := keeper.cdc.Marshal(&brand); err != nil {
		return nil, err
	} else {
		return bz, err
	}
}

// UnmarshalBrand defines a method for returning brand from raw encoded brand
func (keeper Keeper) UnmarshalBrand(bz []byte, brand *types.Brand) error {
	return keeper.cdc.Unmarshal(bz, brand)
}

// SetBrandByOwner defines a method for indexing brand ids by owner
func (k Keeper) SetBrandByOwner(ctx sdk.Context, brandID string, owner sdk.AccAddress) {
	ownerStore := k.getBrandByOwnerStore(ctx, owner)

	ownerStore.Set([]byte(brandID), []byte(brandID))
}

// DeleteBrandByOwner defines a method for removed indexed brand by owner
func (k Keeper) DeleteBrandByOwner(ctx sdk.Context, brandID string, owner sdk.AccAddress) {
	ownerStore := k.getBrandByOwnerStore(ctx, owner)

	ownerStore.Delete([]byte(brandID))
}
