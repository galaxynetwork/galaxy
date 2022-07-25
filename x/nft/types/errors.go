package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidClassID          = sdkerrors.Register(ModuleName, 2, "invalid class id")
	ErrInvalidFeeBasisPoints   = sdkerrors.Register(ModuleName, 3, "invalid fee basis_points")
	ErrInvalidClassDescription = sdkerrors.Register(ModuleName, 4, "invalid class description")
	ErrInvalidNFTID            = sdkerrors.Register(ModuleName, 5, "invalid nft id")
	ErrInvalidNFTUri           = sdkerrors.Register(ModuleName, 6, "invalid nft uri")
	ErrInvalidNFTVarUri        = sdkerrors.Register(ModuleName, 7, "invalid nft variable uri")
	ErrExistClass              = sdkerrors.Register(ModuleName, 8, "exist class within the brand")
	ErrNotFoundClass           = sdkerrors.Register(ModuleName, 9, "not found class within the brand")
	ErrUnauthorized            = sdkerrors.Register(ModuleName, 10, "unauthorized")
	ErrExistNFT                = sdkerrors.Register(ModuleName, 11, "exist nft within the class of brand")
	ErrAlreadyBurnedNFT        = sdkerrors.Register(ModuleName, 12, "already burned nft within the class of brand")
	ErrNotFoundNFT             = sdkerrors.Register(ModuleName, 13, "not found nft within the class of brand")
	ErrNotFoundOwner           = sdkerrors.Register(ModuleName, 14, "not found owner of nft")
)
