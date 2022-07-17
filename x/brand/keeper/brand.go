package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

// HasBrand determines whether the specified classID exist
func (keeper Keeper) HasBrand(ctx sdk.Context, brandID string) bool {
	return false
}

// SetBrand defines a method for set an brand in the store.
func (keeper Keeper) SetBrand(ctx sdk.Context, brand types.Brand) {
}

// GetBrand defines a method for returning a existing brand
func (keeper Keeper) GetBrand(ctx sdk.Context, brandID string) (types.Brand, bool) {
	return types.Brand{}, false
}

func (keeper Keeper) IterateBrands(ctx sdk.Context, cb func(brand types.Brand) (stop bool)) {
}

// GetBrands defines a method for returning all brands
func (keeper Keeper) GetBrands(ctx sdk.Context) (brands types.Brands) {
	return nil
}

// GetBrandsByOwner defines a method for returning all brands by owner
func (keeper Keeper) GetBrandsByOwner(ctx sdk.Context, owner string) types.Brands {
	return nil
}

// MarshalBrand defines a method for protobuf serializes brand
func (keeper Keeper) MarshalBrand(proposal types.Brand) ([]byte, error) {
	return nil, nil
}

// UnmarshalBrand defines a method for returning brand from raw encoded brand
func (keeper Keeper) UnmarshalBrand(bz []byte, brand *types.Brand) error {
	return nil
}
