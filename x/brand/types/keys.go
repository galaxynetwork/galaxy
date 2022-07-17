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

var (
	KeyPrefixBrand = []byte{0x01}

	KeyPrefixBrandByOwner = []byte{0x02}
)

func GetBrandKey(brandID string) []byte {
	return []byte(brandID)
}

func GetPrefixBrandByOwnerKey(owner sdk.AccAddress) []byte {
	return append(KeyPrefixBrandByOwner, address.MustLengthPrefix(owner)...)
}

func NewBrandAddress(brandID string) sdk.AccAddress {
	return address.Module(ModuleName, []byte(brandID))
}
