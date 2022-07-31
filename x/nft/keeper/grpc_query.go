package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/galaxies-labs/galaxy/internal/conv"
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
		if pageRes, err = query.Paginate(k.getClassOfBrandStore(ctx, req.BrandId), req.Pagination, func(_ []byte, bz []byte) error {
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
		return nil, status.Errorf(codes.NotFound, "not found class for brandID: %s, id: %s", req.BrandId, req.ClassId)
	}

	return &types.QueryClassResponse{Class: class}, nil
}

// NFTs queries all nfts belonging to a given brand and class
func (k Querier) NFTs(goCtx context.Context, req *types.QueryNFTsRequest) (*types.QueryNFTsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var owner sdk.AccAddress
	var err error

	if len(req.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner address: %s", req.Owner)
		}
	}

	if len(req.ClassId) > 0 {
		if len(req.BrandId) == 0 {
			return nil, status.Error(codes.InvalidArgument, "order to search by classID, the brandID cannot be empty")
		} else {
			if err = brandtypes.ValidateBrandID(req.BrandId); err != nil {
				return nil, status.Error(codes.InvalidArgument, brandtypes.ErrInvalidBrandID.Error())
			}
		}
		if err = types.ValidateClassId(req.ClassId); err != nil {
			return nil, status.Error(codes.InvalidArgument, types.ErrInvalidClassID.Error())
		}
	} else {
		if len(req.BrandId) > 0 {
			return nil, status.Error(codes.InvalidArgument, "sub classID cannot be empty to search by brandID")
		}
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var pageRes *query.PageResponse
	var nfts types.NFTs

	switch {
	//for paging by key prefix so check owner first
	case len(req.Owner) > 0:
		var store prefix.Store
		if len(req.ClassId) > 0 {
			store = k.getClassStoreByOwner(ctx, owner, req.BrandId, req.ClassId)
		} else {
			store = k.getNFTStoreByOwner(ctx, owner)
		}
		if pageRes, err = query.Paginate(store, req.Pagination, func(key, value []byte) error {
			brandID, classId, id, err := types.ParseNFTUniqueID(conv.UnsafeBytesToStr(value))
			if err != nil {
				panic("invalid nftUniqueID stored in nftOfClassByOwner")
			}
			nft, exist := k.GetNFT(ctx, brandID, classId, id)
			if !exist {
				panic("unexpected nft is stored in nftOfClassByOwner store")
			}

			nfts = append(nfts, nft)
			return nil
		}); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case len(req.ClassId) > 0:
		if pageRes, err = query.Paginate(k.getNFTStore(ctx, req.BrandId, req.ClassId), req.Pagination, func(key, value []byte) error {
			var nft types.NFT
			if err := k.cdc.Unmarshal(value, &nft); err != nil {
				return err
			}

			nfts = append(nfts, nft)
			return nil
		}); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	default:
		if pageRes, err = query.Paginate(k.getNFTPrefix(ctx), req.Pagination, func(key, value []byte) error {
			var nft types.NFT
			if err := k.cdc.Unmarshal(value, &nft); err != nil {
				return err
			}

			nfts = append(nfts, nft)
			return nil
		}); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &types.QueryNFTsResponse{Nfts: nfts, Pagination: pageRes}, nil
}

// NFT queries based on it's brand and class and id
func (k Querier) NFT(goCtx context.Context, req *types.QueryNFTRequest) (*types.QueryNFTResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := brandtypes.ValidateBrandID(req.BrandId); err != nil {
		return nil, status.Error(codes.InvalidArgument, brandtypes.ErrInvalidBrandID.Error())
	}

	if err := types.ValidateClassId(req.ClassId); err != nil {
		return nil, status.Error(codes.InvalidArgument, types.ErrInvalidClassID.Error())
	}

	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, types.ErrInvalidNFTID.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	nft, exist := k.GetNFT(ctx, req.BrandId, req.ClassId, req.Id)
	if !exist {
		return nil, status.Errorf(codes.NotFound, "not found nft for brandID: %s, classID: %s, id: %d", req.BrandId, req.ClassId, req.Id)
	}
	return &types.QueryNFTResponse{Nft: nft}, nil
}

// Owner queries the owner of the NFT based on its brand and class and id
func (k Querier) Owner(goCtx context.Context, req *types.QueryOwnerRequest) (*types.QueryOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := brandtypes.ValidateBrandID(req.BrandId); err != nil {
		return nil, status.Error(codes.InvalidArgument, brandtypes.ErrInvalidBrandID.Error())
	}

	if err := types.ValidateClassId(req.ClassId); err != nil {
		return nil, status.Error(codes.InvalidArgument, types.ErrInvalidClassID.Error())
	}

	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, types.ErrInvalidNFTID.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := k.GetOwner(ctx, req.BrandId, req.ClassId, req.Id)

	if owner == nil {
		return nil, status.Error(codes.NotFound, types.ErrNotFoundClass.Error())
	}

	return &types.QueryOwnerResponse{Owner: owner.String()}, nil
}

// Supply queries the number of NFTs from the given brand and class id
func (k Querier) Supply(goCtx context.Context, req *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := brandtypes.ValidateBrandID(req.BrandId); err != nil {
		return nil, status.Error(codes.InvalidArgument, brandtypes.ErrInvalidBrandID.Error())
	}
	if err := types.ValidateClassId(req.ClassId); err != nil {
		return nil, status.Error(codes.InvalidArgument, types.ErrInvalidClassID.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	supply, err := k.GetSupply(ctx, req.BrandId, req.ClassId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySupplyResponse{NextSequence: supply.Sequence, Amount: supply.TotalSupply}, nil
}
