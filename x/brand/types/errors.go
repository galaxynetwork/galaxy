package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidBrandID           = sdkerrors.Register(ModuleName, 2, "invalid brand id")
	ErrInvalidBrandAddress      = sdkerrors.Register(ModuleName, 3, "invalid brand address")
	ErrInvalidBrandOwnerAddress = sdkerrors.Register(ModuleName, 4, "invalid brand owner address")
	ErrInvalidBrandName         = sdkerrors.Register(ModuleName, 5, "invalid brand name")
	ErrInvalidBrandDetails      = sdkerrors.Register(ModuleName, 6, "invalid brand details")
	ErrInvalidBrandImageUri     = sdkerrors.Register(ModuleName, 7, "invalid brand image_uri")
	ErrExistBrandID             = sdkerrors.Register(ModuleName, 8, "exist brand id")
	ErrExistBrandAddress        = sdkerrors.Register(ModuleName, 9, "exist brand address")
	ErrNotFoundBrand            = sdkerrors.Register(ModuleName, 10, "not found brand")
	ErrUnauthorized             = sdkerrors.Register(ModuleName, 11, "unauthorized")
)
