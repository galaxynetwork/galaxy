package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	ModuleName = "brand"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName

	DefaultBrandCreationFeeDenom = "uglx"
)

var (
	PlaceHolder = []byte{0x01}

	KeyPrefixBrand = []byte{0x01}

	KeyPrefixBrandByOwner = []byte{0x02}
)

func GetPrefixBrandByOwnerKey(owner sdk.AccAddress) []byte {
	return append(KeyPrefixBrandByOwner, address.MustLengthPrefix(owner)...)
}

func NewBrandAddress(brandID string) sdk.AccAddress {
	return address.Module(ModuleName, []byte(brandID))
}
