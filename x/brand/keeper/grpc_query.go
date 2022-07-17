package keeper

import (
	"context"

	"github.com/galaxies-labs/galaxy/x/brand/types"
)

// Brands queries all Brands
func (keeper Keeper) Brands(ctx context.Context, in *types.QueryBrandsRequest) (*types.QueryBrandsResponse, error) {
	return nil, nil
}

// Brand queries and Brand based on it's id
func (keeper Keeper) Brand(ctx context.Context, in *types.QueryBrandRequest) (*types.QueryBrandResponse, error) {
	return nil, nil

}

// BrandsByOwner queries all Brands by owner address
func (keeper Keeper) BrandsByOwner(ctx context.Context, in *types.QueryBrandsByOwnerRequest) (*types.QueryBrandsByOwnerResponse, error) {
	return nil, nil
}
