package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

var _ types.QueryServer = Keeper{}

// Classes queries all Classes
func (k Keeper) Classes(goCtx context.Context, req *types.QueryClassesRequest) (*types.QueryClassesResponse, error) {
	var err error

	if len(req.BrandId) > 0 {
		if err = brandtypes.ValidateBrandID(req.BrandId); err != nil {
			return nil, brandtypes.ErrInvalidBrandID
		}
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var classes types.Classes
	var pageRes *query.PageResponse

	switch {
	case len(req.BrandId) > 0:
		if pageRes, err = query.Paginate(k.getBrandClassesStore(ctx, req.BrandId), req.Pagination, func(_ []byte, bz []byte) error {
			var class types.Class
			if err := k.cdc.Unmarshal(bz, &class); err != nil {
				return err
			}

			classes = append(classes, class)
			return nil
		}); err != nil {
			return nil, err
		}
	default:
		if pageRes, err = query.Paginate(k.getClassesStore(ctx), req.Pagination, func(_ []byte, bz []byte) error {
			var class types.Class
			if err := k.cdc.Unmarshal(bz, &class); err != nil {
				return err
			}

			classes = append(classes, class)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return &types.QueryClassesResponse{Classes: classes, Pagination: pageRes}, nil
}

// Class queries based on it's id
func (k Keeper) Class(goCtx context.Context, req *types.QueryClassRequest) (*types.QueryClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	class, exist := k.GetClass(ctx, req.ClassId, req.BrandId)
	if !exist {
		return nil, types.ErrNotFoundClass
	}

	return &types.QueryClassResponse{Class: class}, nil
}
