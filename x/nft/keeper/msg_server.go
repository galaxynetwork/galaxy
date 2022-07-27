package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
	class := types.NewClass(msg.BrandId, msg.Id, msg.FeeBasisPoints, msg.Description)
	if err := class.Validate(); err != nil {
		return nil, err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid creater address: %s", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	brand, brandExist := k.brandKeeper.GetBrand(ctx, msg.BrandId)
	if !brandExist {
		return nil, brandtypes.ErrNotFoundBrand
	}

	if brand.Owner != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	if err := k.SaveClass(ctx, class); err != nil {
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
	//validation
	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return nil, err
	}

	if err := types.ValidateClassId(msg.Id); err != nil {
		return nil, err
	}

	if err := types.ValidateFeeBasisPoints(msg.FeeBasisPoints); err != nil {
		return nil, err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Editor); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid editor address: %s", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	brand, exist := k.brandKeeper.GetBrand(ctx, msg.BrandId)
	if !exist {
		return nil, brandtypes.ErrNotFoundBrand
	}

	if brand.Owner != msg.Editor {
		return nil, types.ErrUnauthorized
	}

	class, exist := k.GetClass(ctx, msg.BrandId, msg.Id)
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

// MintNFT defines a method for minting a new nft.
func (k msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	return nil, nil
}

// UpdateNFT defines a method for updating variableURI an existing nft.
func (k msgServer) UpdateNFT(goCtx context.Context, msg *types.MsgUpdateNFT) (*types.MsgUpdateNFTResponse, error) {
	return nil, nil
}

// BurnNFT defines a method for burnning an existing nft.
func (k msgServer) BurnNFT(goCtx context.Context, msg *types.MsgBurnNFT) (*types.MsgBurnNFTResponse, error) {
	return nil, nil
}

// TransferNFT defines a method for transferring ownership an existing nft.
func (k msgServer) TransferNFT(goCtx context.Context, msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
	return nil, nil
}
