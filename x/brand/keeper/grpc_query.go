package keeper

import (
	"context"
	"fmt"

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
func (k Querier) Brands(ctx context.Context, req *types.QueryBrandsRequest) (*types.QueryBrandsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var err error
	var addr sdk.AccAddress

	if len(req.Owner) > 0 {
		addr, err = sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner address: %s", err)
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var brands types.Brands
	var pageRes *query.PageResponse

	switch {
	case len(addr) > 0:
		if pageRes, err = query.Paginate(k.getBrandByOwnerStore(sdkCtx, addr), req.Pagination, func(key, _ []byte) error {
			brand, exist := k.GetBrand(sdkCtx, string(key[:]))
			if !exist {
				return fmt.Errorf("unexpected brand is stored in brand_by_owner store")
			}

			brands = append(brands, brand)

			return nil
		}); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	default:
		if pageRes, err = query.Paginate(k.getBrandStore(sdkCtx), req.Pagination, func(_, bz []byte) error {
			var brand types.Brand

			if err := k.UnmarshalBrand(bz, &brand); err != nil {
				return err
			}

			brands = append(brands, brand)

			return nil
		}); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	return &types.QueryBrandsResponse{Brands: brands, Pagination: pageRes}, nil
}

// Brand queries and Brand based on it's id
func (keeper Querier) Brand(ctx context.Context, req *types.QueryBrandRequest) (*types.QueryBrandResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := types.ValidateBrandID(req.BrandId); err != nil {
		return nil, status.Error(codes.InvalidArgument, types.ErrInvalidBrandID.Error())
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	brand, _ := keeper.GetBrand(sdkCtx, req.BrandId)

	return &types.QueryBrandResponse{Brand: brand}, nil

}
