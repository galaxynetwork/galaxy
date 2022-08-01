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

	if err := types.ValidateFeeBasisPoints(msg.FeeBasisPoints, true); err != nil {
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

	class = class.UpdateFeeBasisPoints(msg.FeeBasisPoints)
	class.Description = desc

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

// MintToNFT defines a method for minting a new nft.
func (k msgServer) MintToNFT(goCtx context.Context, msg *types.MsgMintToNFT) (*types.MsgMintToNFTResponse, error) {
	nft := types.NewNFT(1, msg.BrandId, msg.ClassId, msg.Uri, msg.VarUri)

	//for basic validation
	if err := nft.Validate(); err != nil {
		return nil, err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Minter); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid minter address: %s", err)
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid recipient address: %s", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if brand, exist := k.brandKeeper.GetBrand(ctx, msg.BrandId); !exist {
		return nil, brandtypes.ErrNotFoundBrand
	} else {
		if brand.Owner != msg.Minter {
			return nil, types.ErrUnauthorized
		}
	}

	if exist := k.HasClass(ctx, msg.BrandId, msg.ClassId); !exist {
		return nil, types.ErrNotFoundClass
	}

	if sequence, err := k.getSequenceOfClass(ctx, msg.BrandId, msg.ClassId); err != nil {
		return nil, err
	} else {
		nft.Id = sequence
	}

	if err := k.Keeper.MintToNFT(ctx, nft, recipient); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgMintToNFT,
			sdk.NewAttribute(brandtypes.AttributeBrandID, nft.BrandId),
			sdk.NewAttribute(types.AttributeClassID, nft.ClassId),
			sdk.NewAttribute(types.AttributeNFTID, strconv.FormatUint(nft.Id, 10)),
			sdk.NewAttribute(types.AttributeRecipient, msg.Recipient),
		),
	)

	return &types.MsgMintToNFTResponse{Id: nft.Id}, nil
}

// UpdateNFT defines a method for updating variableURI an existing nft.
func (k msgServer) UpdateNFT(goCtx context.Context, msg *types.MsgUpdateNFT) (*types.MsgUpdateNFTResponse, error) {
	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return nil, err
	}

	if err := types.ValidateClassId(msg.ClassId); err != nil {
		return nil, err
	}

	if msg.Id == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNFTID, "nft id must be greater 0")
	}

	if len(msg.VarUri) > types.MaxUriLength {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNFTVarUri, "invalid var_uri length; got: %d, max: %d", len(msg.VarUri), types.MaxUriLength)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	owenr := k.GetOwner(ctx, msg.BrandId, msg.ClassId, msg.Id)
	if owenr == nil {
		return nil, types.ErrNotFoundNFT
	}

	if msg.Sender != owenr.String() {
		return nil, types.ErrUnauthorized
	}

	if err := k.Keeper.UpdateNFT(ctx, msg.BrandId, msg.ClassId, msg.Id, msg.VarUri); err != nil {
		return nil, err
	}

	return &types.MsgUpdateNFTResponse{}, nil
}

// BurnNFT defines a method for burnning an existing nft.
func (k msgServer) BurnNFT(goCtx context.Context, msg *types.MsgBurnNFT) (*types.MsgBurnNFTResponse, error) {
	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return nil, err
	}

	if err := types.ValidateClassId(msg.ClassId); err != nil {
		return nil, err
	}

	if msg.Id == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNFTID, "nft id must be greater 0")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	owenr := k.GetOwner(ctx, msg.BrandId, msg.ClassId, msg.Id)
	if owenr == nil {
		return nil, types.ErrNotFoundNFT
	}

	if msg.Sender != owenr.String() {
		return nil, types.ErrUnauthorized
	}

	if err := k.Keeper.BurnNFT(ctx, msg.BrandId, msg.ClassId, msg.Id); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgBurnNFT,
			sdk.NewAttribute(brandtypes.AttributeBrandID, msg.BrandId),
			sdk.NewAttribute(types.AttributeClassID, msg.ClassId),
			sdk.NewAttribute(types.AttributeNFTID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeSender, msg.Sender),
		),
	)

	return &types.MsgBurnNFTResponse{}, nil
}

// TransferNFT defines a method for transferring ownership an existing nft.
func (k msgServer) TransferNFT(goCtx context.Context, msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return nil, err
	}

	if err := types.ValidateClassId(msg.ClassId); err != nil {
		return nil, err
	}

	if msg.Id == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidNFTID, "nft id must be greater 0")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid recipient address: %s", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	owenr := k.GetOwner(ctx, msg.BrandId, msg.ClassId, msg.Id)
	if owenr == nil {
		return nil, types.ErrNotFoundNFT
	}

	if msg.Sender != owenr.String() {
		return nil, types.ErrUnauthorized
	}

	if err := k.Keeper.TransferNFT(ctx, msg.BrandId, msg.ClassId, msg.Id, owenr, recipient); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgTransferNFT,
			sdk.NewAttribute(brandtypes.AttributeBrandID, msg.BrandId),
			sdk.NewAttribute(types.AttributeClassID, msg.ClassId),
			sdk.NewAttribute(types.AttributeNFTID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeSender, msg.Sender),
			sdk.NewAttribute(types.AttributeRecipient, msg.Recipient),
		),
	)

	return &types.MsgTransferNFTResponse{}, nil
}
