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
)
