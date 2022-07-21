package types

const (
	ModuleName = "nft"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

var (
	PrefixClassKey = []byte{0x01}
)

func GetClassStoreKey(brandId string) []byte {
	return append(PrefixClassKey, []byte(brandId)...)
}
