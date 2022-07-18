package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreateBrand{}
	_ sdk.Msg = &MsgEditBrand{}
	_ sdk.Msg = &MsgTransferOwnershipBrand{}
)

const (
	TypeMsgCreateBrand            = "create_brand"
	TypeMsgEditBrand              = "edit_brand"
	TypeMsgTransferOwnershipBrand = "transfer_ownership_brand"
)

func NewMsgCreateBrand(id, owner string, description BrandDescription) *MsgCreateBrand {
	return &MsgCreateBrand{
		Id:          id,
		Owner:       owner,
		Description: description,
	}
}

func (msg MsgCreateBrand) Route() string { return RouterKey }

func (msg MsgCreateBrand) Type() string { return TypeMsgCreateBrand }

func (msg MsgCreateBrand) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgCreateBrand) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}

	if err := ValidateBrandID(msg.Id); err != nil {
		return sdkerrors.Wrapf(ErrInvalidBrandID, "invalid brand id: %s", err)
	}

	if err := msg.Description.Validate(); err != nil {
		return sdkerrors.Wrap(err, "invalid brand description")
	}

	return nil
}

func (msg MsgCreateBrand) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	return []sdk.AccAddress{owner}
}

func NewMsgEditBrand(id, owner string, description BrandDescription) *MsgEditBrand {
	return &MsgEditBrand{
		Id:          id,
		Owner:       owner,
		Description: description,
	}
}

func (msg MsgEditBrand) Route() string { return RouterKey }

func (msg MsgEditBrand) Type() string { return TypeMsgEditBrand }

func (msg MsgEditBrand) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgEditBrand) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}

	if err := ValidateBrandID(msg.Id); err != nil {
		return sdkerrors.Wrapf(ErrInvalidBrandID, "invalid brand id: %s", err)
	}

	if err := msg.Description.Validate(); err != nil {
		return sdkerrors.Wrap(err, "invalid brand description")
	}

	return nil
}

func (msg MsgEditBrand) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	return []sdk.AccAddress{owner}
}

func NewMsgTransferOwnershipBrand(id, owner, destOwner string) *MsgTransferOwnershipBrand {
	return &MsgTransferOwnershipBrand{
		Id:        id,
		Owner:     owner,
		DestOwner: destOwner,
	}
}

func (msg MsgTransferOwnershipBrand) Route() string { return RouterKey }

func (msg MsgTransferOwnershipBrand) Type() string { return TypeMsgTransferOwnershipBrand }

func (msg MsgTransferOwnershipBrand) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(&msg)
}

func (msg MsgTransferOwnershipBrand) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.DestOwner); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid dest owner address: %s", err)
	}

	if err := ValidateBrandID(msg.Id); err != nil {
		return sdkerrors.Wrapf(ErrInvalidBrandID, "invalid brand id: %s", err)
	}

	return nil
}

func (msg MsgTransferOwnershipBrand) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	return []sdk.AccAddress{owner}
}
