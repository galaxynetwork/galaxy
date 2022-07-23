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
	KeyPrefixBrand        = []byte{0x01}
	KeyPrefixBrandByOwner = []byte{0x02}

	Delimiter   = []byte{0x00}
	PlaceHolder = []byte{0x01}
)

func GetPrefixBrandByOwnerKey(owner sdk.AccAddress) []byte {
	ownerLp := address.MustLengthPrefix(owner)
	key := make([]byte, len(KeyPrefixBrandByOwner)+len(ownerLp)+len(Delimiter))
	copy(key, KeyPrefixBrandByOwner)
	copy(key[len(KeyPrefixBrandByOwner):], ownerLp)
	copy(key[len(KeyPrefixBrandByOwner)+len(ownerLp):], Delimiter)
	return key
}

func NewBrandAddress(brandID string) sdk.AccAddress {
	return address.Module(ModuleName, []byte(brandID))
}
