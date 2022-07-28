package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	brandtypes "github.com/galaxies-labs/galaxy/x/brand/types"
)

var (
	_ sdk.Msg = &MsgCreateClass{}
	_ sdk.Msg = &MsgEditClass{}

	_ sdk.Msg = &MsgMintNFT{}
	_ sdk.Msg = &MsgBurnNFT{}
	_ sdk.Msg = &MsgUpdateNFT{}
	_ sdk.Msg = &MsgTransferNFT{}
)

const (
	TypeMsgCreateClass = "crete-class"
	TypeMsgEditClass   = "edit-class"

	TypeMsgMintNFT     = "mint-nft"
	TypeMsgBurnNFT     = "burn-nft"
	TypeMsgUpdateNFT   = "update-nft"
	TypeMsgTransferNFT = "transfer-nft"
)

func NewMsgCreateClass(brandID, id, creator string, feeBasisPoints uint32, description ClassDescription) *MsgCreateClass {
	return &MsgCreateClass{
		BrandId:        brandID,
		Id:             id,
		FeeBasisPoints: feeBasisPoints,
		Description:    description,
		Creator:        creator,
	}
}

func (msg MsgCreateClass) Route() string { return RouterKey }

func (msg MsgCreateClass) Type() string { return TypeMsgCreateClass }

func (msg MsgCreateClass) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgCreateClass) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid creator address: %s", err)
	}

	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return err
	}

	if err := ValidateClassId(msg.Id); err != nil {
		return err
	}

	if err := msg.Description.Validate(); err != nil {
		return err
	}

	if msg.FeeBasisPoints == 0 || msg.FeeBasisPoints > MaxFeeBasisPoints {
		return sdkerrors.Wrapf(ErrInvalidFeeBasisPoints, " got: %d, min: %d, max: %d", msg.FeeBasisPoints, 1, MaxFeeBasisPoints)
	}

	return nil
}

func (msg MsgCreateClass) GetSigners() []sdk.AccAddress {
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func NewMsgEditClass(brandID, id, editor string, feeBasisPoints uint32, description ClassDescription) *MsgEditClass {
	return &MsgEditClass{
		BrandId:        brandID,
		Id:             id,
		FeeBasisPoints: feeBasisPoints,
		Description:    description,
		Editor:         editor,
	}
}

func (msg MsgEditClass) Route() string { return RouterKey }

func (msg MsgEditClass) Type() string { return TypeMsgEditClass }

func (msg MsgEditClass) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgEditClass) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Editor); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid editor address: %s", err)
	}

	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return err
	}

	if err := ValidateClassId(msg.Id); err != nil {
		return err
	}

	if err := msg.Description.Validate(); err != nil {
		return err
	}

	if msg.FeeBasisPoints == 0 || (msg.FeeBasisPoints != DoNotModifyFeeBasisPoints && msg.FeeBasisPoints > MaxFeeBasisPoints) {
		return sdkerrors.Wrapf(ErrInvalidFeeBasisPoints, " got: %d, min: %d, max: %d", msg.FeeBasisPoints, 1, MaxFeeBasisPoints)
	}

	return nil
}

func (msg MsgEditClass) GetSigners() []sdk.AccAddress {
	editor, _ := sdk.AccAddressFromBech32(msg.Editor)
	return []sdk.AccAddress{editor}
}

func NewMsgMintNFT(brandID, classId, uri, varUri, minter, recipient string) *MsgMintNFT {
	return &MsgMintNFT{
		BrandId:   brandID,
		ClassId:   classId,
		Uri:       uri,
		VarUri:    varUri,
		Minter:    minter,
		Recipient: recipient,
	}
}

func (msg MsgMintNFT) Route() string { return RouterKey }
func (msg MsgMintNFT) Type() string  { return TypeMsgMintNFT }
func (msg MsgMintNFT) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgMintNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Minter); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid minter address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid recipient address: %s", err)
	}

	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return err
	}

	if err := ValidateClassId(msg.ClassId); err != nil {
		return err
	}

	if len(msg.Uri) == 0 {
		return sdkerrors.Wrap(ErrInvalidNFTUri, "uri can not be blank")
	}

	if len(msg.Uri) > MaxUriLength {
		return sdkerrors.Wrapf(ErrInvalidNFTUri, "invalid uri length; got: %d, max: %d", len(msg.Uri), MaxUriLength)
	}

	if len(msg.VarUri) > MaxUriLength {
		return sdkerrors.Wrapf(ErrInvalidNFTVarUri, "invalid var_uri length; got: %d, max: %d", len(msg.VarUri), MaxUriLength)
	}

	return nil
}

func (msg MsgMintNFT) GetSigners() []sdk.AccAddress {
	minter, _ := sdk.AccAddressFromBech32(msg.Minter)
	return []sdk.AccAddress{minter}
}

func NewMsgBurnNFT(brandID, classId string, id uint64, sender string) *MsgBurnNFT {
	return &MsgBurnNFT{
		BrandId: brandID,
		ClassId: classId,
		Id:      id,
		Sender:  sender,
	}
}

func (msg MsgBurnNFT) Route() string { return RouterKey }
func (msg MsgBurnNFT) Type() string  { return TypeMsgBurnNFT }
func (msg MsgBurnNFT) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgBurnNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return err
	}

	if err := ValidateClassId(msg.ClassId); err != nil {
		return err
	}

	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidNFTID, "nft id must be greater 0")
	}

	return nil
}

func (msg MsgBurnNFT) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

func NewMsgTransferNFT(brandID, classId string, id uint64, sender, recipient string) *MsgTransferNFT {
	return &MsgTransferNFT{
		BrandId:   brandID,
		ClassId:   classId,
		Id:        id,
		Sender:    sender,
		Recipient: recipient,
	}
}

func (msg MsgTransferNFT) Route() string { return RouterKey }
func (msg MsgTransferNFT) Type() string  { return TypeMsgTransferNFT }
func (msg MsgTransferNFT) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgTransferNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid minter address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid recipient address: %s", err)
	}

	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return err
	}

	if err := ValidateClassId(msg.ClassId); err != nil {
		return err
	}

	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidNFTID, "nft id must be greater 0")
	}

	return nil
}

func (msg MsgTransferNFT) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

func NewMsgUpdateNFT(brandID, classId string, id uint64, varUri, sender string) *MsgUpdateNFT {
	return &MsgUpdateNFT{
		BrandId: brandID,
		ClassId: classId,
		Id:      id,
		VarUri:  varUri,
		Sender:  sender,
	}
}

func (msg MsgUpdateNFT) Route() string { return RouterKey }
func (msg MsgUpdateNFT) Type() string  { return TypeMsgUpdateNFT }
func (msg MsgUpdateNFT) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgUpdateNFT) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}

	if err := brandtypes.ValidateBrandID(msg.BrandId); err != nil {
		return err
	}

	if err := ValidateClassId(msg.ClassId); err != nil {
		return err
	}

	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalidNFTID, "nft id must be greater 0")
	}

	if len(msg.VarUri) > MaxUriLength {
		return sdkerrors.Wrapf(ErrInvalidNFTVarUri, "invalid var_uri length; got: %d, max: %d", len(msg.VarUri), MaxUriLength)
	}

	return nil
}

func (msg MsgUpdateNFT) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}
