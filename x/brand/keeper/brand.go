package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

// HasBrand determines whether the specified classID exist
func (k Keeper) HasBrand(ctx sdk.Context, brandID string) bool {
	brandStore := k.getBrandStore(ctx)

	return brandStore.Has([]byte(brandID))
}

// SetBrand defines a method for set an brand in the store.
func (k Keeper) SetBrand(ctx sdk.Context, brand types.Brand) error {
	brandStore := k.getBrandStore(ctx)

	bz, err := k.MarshalBrand(brand)
	if err != nil {
		return err
	}

	brandStore.Set([]byte(brand.Id), bz)

	return nil
}

// GetBrand defines a method for returning a existing brand
func (k Keeper) GetBrand(ctx sdk.Context, brandID string) (types.Brand, bool) {
	brandStore := k.getBrandStore(ctx)

	var brand types.Brand

	bz := brandStore.Get([]byte(brandID))
	if bz == nil {
		return brand, false
	}

	if err := k.UnmarshalBrand(bz, &brand); err != nil {
		panic(fmt.Errorf("stored brand unmarshalling error: %v", err))
	}

	return brand, true
}

func (k Keeper) IterateBrands(ctx sdk.Context, cb func(brand types.Brand) (stop bool)) {
	brandStore := k.getBrandStore(ctx)

	iterator := brandStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var brand types.Brand
		if err := k.UnmarshalBrand(iterator.Value(), &brand); err != nil {
			panic(fmt.Errorf("stored brand unmarshalling error: %v", err))
		}

		if cb(brand) {
			break
		}
	}
}

func (k Keeper) IterateBrandsByOwner(ctx sdk.Context, owner sdk.AccAddress, cb func(brand types.Brand) (stop bool)) {
	ownerStore := k.getBrandByOwnerStore(ctx, owner)

	iterator := ownerStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		brand, exist := k.GetBrand(ctx, string(iterator.Key()[:]))
		if !exist {
			panic("unexpected brand is stored in brand_by_owner store")
		}

		if cb(brand) {
			break
		}
	}
}

// GetBrands defines a method for returning all brands
func (k Keeper) GetBrands(ctx sdk.Context) (brands types.Brands) {
	k.IterateBrands(ctx, func(brand types.Brand) (stop bool) {
		brands = append(brands, brand)
		return false
	})
	return
}

// GetBrands defines a method for returning all brands a given owner
func (k Keeper) GetBrandsByOwner(ctx sdk.Context, owner sdk.AccAddress) (brands types.Brands) {
	k.IterateBrandsByOwner(ctx, owner, func(brand types.Brand) (stop bool) {
		brands = append(brands, brand)
		return false
	})
	return
}

// MarshalBrand defines a method for protobuf serializes brand
func (k Keeper) MarshalBrand(brand types.Brand) ([]byte, error) {
	if bz, err := k.cdc.Marshal(&brand); err != nil {
		return nil, err
	} else {
		return bz, err
	}
}

// UnmarshalBrand defines a method for returning brand from raw encoded brand
func (k Keeper) UnmarshalBrand(bz []byte, brand *types.Brand) error {
	return k.cdc.Unmarshal(bz, brand)
}

// SetBrandByOwner defines a method for indexing brand ids by owner
func (k Keeper) SetBrandByOwner(ctx sdk.Context, brandID string, owner sdk.AccAddress) {
	ownerStore := k.getBrandByOwnerStore(ctx, owner)

	ownerStore.Set([]byte(brandID), types.PlaceHolder)
}

// DeleteBrandByOwner defines a method for removed indexed brand by owner
func (k Keeper) DeleteBrandByOwner(ctx sdk.Context, brandID string, owner sdk.AccAddress) {
	ownerStore := k.getBrandByOwnerStore(ctx, owner)

	ownerStore.Delete([]byte(brandID))
}

// getBrandStore get brand store
func (k Keeper) getBrandStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixBrand)
}

// getBrandByOwnerStore get owner specific brand store
func (k Keeper) getBrandByOwnerStore(ctx sdk.Context, owner sdk.AccAddress) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.GetPrefixBrandByOwnerKey(owner))
}
