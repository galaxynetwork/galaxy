package keeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/galaxies-labs/galaxy/x/brand/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// Brands queries all Brands
func (keeper Querier) Brands(ctx context.Context, req *types.QueryBrandsRequest) (*types.QueryBrandsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	brandStore := prefix.NewStore(sdkCtx.KVStore(keeper.storeKey), types.KeyPrefixBrand)

	var brands types.Brands

	pageRes, err := query.Paginate(brandStore, req.Pagination, func(key, value []byte) error {
		var brand types.Brand
		if err := keeper.UnmarshalBrand(value, &brand); err != nil {
			return err
		}

		brands = append(brands, brand)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryBrandsResponse{Brands: brands, Pagination: pageRes}, nil
}

// Brand queries and Brand based on it's id
func (keeper Querier) Brand(ctx context.Context, req *types.QueryBrandRequest) (*types.QueryBrandResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if len(strings.TrimSpace(req.BrandId)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "brand id cannot be empty")
	}

	if err := types.ValidateBrandID(req.BrandId); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid brand id: %s", err.Error())
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	brand, _ := keeper.GetBrand(sdkCtx, req.BrandId)

	return &types.QueryBrandResponse{Brand: brand}, nil

}

// BrandsByOwner queries all Brands by owner address
func (keeper Querier) BrandsByOwner(ctx context.Context, req *types.QueryBrandsByOwnerRequest) (*types.QueryBrandsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if len(strings.TrimSpace(req.Owner)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "owner address cannot be empty")
	}

	addr, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner address: %s", err.Error())
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	ownerStore := keeper.getBrandByOwnerStore(sdkCtx, addr)

	var brands types.Brands

	pageRes, err := query.Paginate(ownerStore, req.Pagination, func(key, value []byte) error {
		brand, exist := keeper.GetBrand(sdkCtx, string(value[:]))
		if !exist {
			panic(fmt.Errorf("unexpected brand is stored in brand_by_owner prefix"))
		}
		brands = append(brands, brand)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryBrandsByOwnerResponse{Brands: brands, Pagination: pageRes}, nil
}
