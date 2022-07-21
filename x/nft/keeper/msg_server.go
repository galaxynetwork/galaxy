package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateClass defines a method for creating a new class within brand.
func (k msgServer) CreateClass(goCtx context.Context, msg *types.MsgCreateClass) (*types.MsgCreateClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	brand, exist := k.brandKeeper.GetBrand(ctx, msg.BrandId)
	if !exist {
		return nil, brandtypes.ErrNotFoundBrand
	}

	if brand.Owner != msg.Creator {
		return nil, brandtypes.ErrUnauthorized
	}

	class := types.NewClass(msg.BrandId, msg.Id, msg.FeeBasisPoints, msg.MaxSupply, msg.Description)
	if err := class.Validate(); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgCreateClass,
			sdk.NewAttribute(brandtypes.AttributeBrandID, class.BrandId),
			sdk.NewAttribute(types.AttributeClassID, class.Id),
		),
	)

	return &types.MsgCreateClassResponse{}, nil
}

// EditClass defines a method for editing an existing class.
func (k msgServer) EditClass(goCtx context.Context, msg *types.MsgEditClass) (*types.MsgEditClassResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	brand, exist := k.brandKeeper.GetBrand(ctx, msg.BrandId)
	if !exist {
		return nil, brandtypes.ErrNotFoundBrand
	}

	if brand.Owner != msg.Editor {
		return nil, brandtypes.ErrUnauthorized
	}

	class, exist := k.GetClass(ctx, msg.Id, msg.BrandId)
	if !exist {
		return nil, types.ErrNotFoundClass
	}

	desc := class.Description.UpdateDescription(msg.Description)
	if err := desc.Validate(); err != nil {
		return nil, err
	}

	class.Description = desc
	class.FeeBasisPoints = msg.FeeBasisPoints

	if err := k.SetClass(ctx, class); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgEditClass,
			sdk.NewAttribute(brandtypes.AttributeBrandID, class.BrandId),
			sdk.NewAttribute(types.AttributeClassID, class.Id),
			sdk.NewAttribute(
				types.AttributeFeeBasisPoints,
				strconv.FormatUint(uint64(class.FeeBasisPoints), 10),
			),
		),
	)

	return &types.MsgEditClassResponse{}, nil

}
