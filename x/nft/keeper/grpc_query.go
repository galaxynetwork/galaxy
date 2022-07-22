package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
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

// Classes queries all Classes
func (k Querier) Classes(goCtx context.Context, req *types.QueryClassesRequest) (*types.QueryClassesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var err error
	if len(req.BrandId) > 0 {
		if err = brandtypes.ValidateBrandID(req.BrandId); err != nil {
			return nil, status.Error(codes.InvalidArgument, brandtypes.ErrInvalidBrandID.Error())
		}
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var classes types.Classes
	var pageRes *query.PageResponse

	switch {
	case len(req.BrandId) > 0:
		if pageRes, err = query.Paginate(k.getClassOfBrandPrefixStore(ctx, req.BrandId), req.Pagination, func(_ []byte, bz []byte) error {
			var class types.Class
			if err := k.cdc.Unmarshal(bz, &class); err != nil {
				return err
			}

			classes = append(classes, class)
			return nil
		}); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	default:
		if pageRes, err = query.Paginate(k.getClassPrefixStore(ctx), req.Pagination, func(_ []byte, bz []byte) error {
			var class types.Class
			if err := k.cdc.Unmarshal(bz, &class); err != nil {
				return err
			}

			classes = append(classes, class)
			return nil
		}); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &types.QueryClassesResponse{Classes: classes, Pagination: pageRes}, nil
}

// Class queries based on it's id
func (k Querier) Class(goCtx context.Context, req *types.QueryClassRequest) (*types.QueryClassResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := brandtypes.ValidateBrandID(req.BrandId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := types.ValidateClassId(req.ClassId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	class, exist := k.GetClass(ctx, req.BrandId, req.ClassId)
	if !exist {
		return nil, status.Errorf(codes.NotFound,
			sdkerrors.Wrapf(types.ErrNotFoundClass, "not found class for brandID: %s, id: %s", req.BrandId, req.ClassId).Error(),
		)
	}

	return &types.QueryClassResponse{Class: class}, nil
}
