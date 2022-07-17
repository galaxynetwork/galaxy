package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	ModuleName = "brand"

	StoreKey = ModuleName

	DefaultBrandCreationFeeDenom = "uglx"
)

func NewBrandAddress(brandID string) sdk.AccAddress {
	return address.Module(ModuleName, []byte(brandID))
}
