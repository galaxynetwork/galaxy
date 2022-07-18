package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/galaxies-labs/galaxy/x/brand/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the x/brand MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateBrand defines a method for creating a new Brand.
func (ms msgServer) CreateBrand(goCtx context.Context, msg *types.MsgCreateBrand) (*types.MsgCreateBrandResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	brand := types.NewBrand(msg.Id, owner, msg.Description)
	if err := brand.Validate(); err != nil {
		return nil, err
	}

	if exist := ms.HasBrand(ctx, brand.Id); exist {
		return nil, types.ErrExistBrandID
	}

	brandAddress, _ := sdk.AccAddressFromBech32(brand.BrandAddress)
	brandAcct := ms.authKeeper.GetAccount(ctx, brandAddress)
	if brandAcct != nil {
		return nil, sdkerrors.Wrap(types.ErrExistBrandAddress, brandAcct.GetAddress().String())
	}

	ms.authKeeper.SetAccount(ctx, ms.authKeeper.NewAccountWithAddress(ctx, brandAddress))

	params := ms.Keeper.GetParams(ctx)
	if params.BrandCreationFee.Amount.IsPositive() {
		if err := ms.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(params.BrandCreationFee), owner); err != nil {
			return nil, err
		}
	}

	ms.SetBrand(ctx, brand)
	ms.SetBrandByOwner(ctx, brand.Id, owner)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgCreateBrand,
			sdk.NewAttribute(types.AttributeBrandID, brand.Id),
			sdk.NewAttribute(types.AttributeBrandAddress, brand.BrandAddress),
			sdk.NewAttribute(types.AttributeOwner, brand.Owner),
		),
	)

	return &types.MsgCreateBrandResponse{BrandAddress: brand.BrandAddress}, nil
}

// EditBrand defines a method for editing an existing brand.
func (ms msgServer) EditBrand(goCtx context.Context, msg *types.MsgEditBrand) (*types.MsgEditBrandResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := types.ValidateBrandID(msg.Id); err != nil {
		return nil, types.ErrExistBrandID
	}

	if err := msg.Description.Validate(); err != nil {
		return nil, types.ErrExistBrandID
	}

	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	brand, exist := ms.GetBrand(ctx, msg.Id)
	if !exist {
		return nil, types.ErrNotFoundBrand
	}

	if brand.Owner != msg.Owner {
		return nil, types.ErrUnauthorized
	}

	brand.Description = msg.Description

	ms.SetBrand(ctx, brand)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgEditBrand,
			sdk.NewAttribute(types.AttributeBrandID, brand.Id),
			sdk.NewAttribute(types.AttributeBrandAddress, brand.BrandAddress),
		),
	)

	return &types.MsgEditBrandResponse{}, nil
}

// TransferOwnershipBrand defines a method for transfer ownership of existing brand
func (ms msgServer) TransferOwnershipBrand(goCtx context.Context, msg *types.MsgTransferOwnershipBrand) (*types.MsgTransferOwnershipBrandResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := types.ValidateBrandID(msg.Id); err != nil {
		return nil, types.ErrExistBrandID
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	destOwner, err := sdk.AccAddressFromBech32(msg.DestOwner)
	if err != nil {
		return nil, err
	}

	brand, exist := ms.GetBrand(ctx, msg.Id)
	if !exist {
		return nil, types.ErrNotFoundBrand
	}

	if brand.Owner != msg.Owner {
		return nil, types.ErrUnauthorized
	}

	brand.Owner = msg.DestOwner

	ms.SetBrand(ctx, brand)
	ms.DeleteBrandByOwner(ctx, brand.Id, owner)
	ms.SetBrandByOwner(ctx, brand.Id, destOwner)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgTransferOwnershipBrand,
			sdk.NewAttribute(types.AttributeBrandID, brand.Id),
			sdk.NewAttribute(types.AttributeBrandAddress, brand.BrandAddress),
			sdk.NewAttribute(types.AttributeNewOwner, brand.Owner),
		),
	)

	return &types.MsgTransferOwnershipBrandResponse{}, nil
}
