package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateClass{}
	_ sdk.Msg = &MsgEditClass{}
)

const (
	TypeMsgCreateClass = "crete-class"
	TypeMsgEditClass   = "edit-class"
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

	return nil
}

func (msg MsgEditClass) GetSigners() []sdk.AccAddress {
	editor, _ := sdk.AccAddressFromBech32(msg.Editor)
	return []sdk.AccAddress{editor}
}
