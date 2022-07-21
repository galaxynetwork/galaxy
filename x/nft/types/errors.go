package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidClassID          = sdkerrors.Register(ModuleName, 2, "invalid class id")
	ErrInvalidFeeBasisPoints   = sdkerrors.Register(ModuleName, 3, "invalid fee basis_points")
	ErrInvalidClassName        = sdkerrors.Register(ModuleName, 4, "invalid class name")
	ErrInvalidClassDetails     = sdkerrors.Register(ModuleName, 5, "invalid class details")
	ErrInvalidClassExternalUrl = sdkerrors.Register(ModuleName, 6, "invalid class external_url")
	ErrInvalidClassImageUri    = sdkerrors.Register(ModuleName, 7, "invalid class image_uri")
	ErrInvalidNFTUri           = sdkerrors.Register(ModuleName, 8, "invalid nft uri")
	ErrInvalidNFTVarUri        = sdkerrors.Register(ModuleName, 9, "invalid nft var_uri")
	ErrExistClassWithinBrand   = sdkerrors.Register(ModuleName, 10, "exist class within the brand")
	ErrNotFoundClass           = sdkerrors.Register(ModuleName, 11, "not found class within the brand")
)
