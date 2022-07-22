package types

import (
	"strings"

	"github.com/galaxies-labs/galaxy/internal/conv"
)

const (
	ModuleName = "nft"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

//(0x01|[]byte(brandID)|0x00|[]byte(classID) -> Class{})
// Query
//	- (0x01) prefix -> All classes
//	- (0x01|[]byte(brandID)|0x00) prefix -> All classes given brandID

var (
	BrandClassKey = []byte{0x01}
	Delimiter     = []byte{0x00}
)

func GetClassUniqueID(brandID, id string) string {
	return strings.Join([]string{brandID, id}, "/")
}

func GetClassOfBrandStoreKey(brandID string) []byte {
	brandIDBz := conv.UnsafeStrToBytes(brandID)

	key := make([]byte, len(BrandClassKey)+len(brandIDBz)+len(Delimiter))

	copy(key, BrandClassKey)
	copy(key[len(BrandClassKey):], brandIDBz)
	copy(key[len(BrandClassKey)+len(brandIDBz):], Delimiter)
	return key
}

func GetPrefixClassKey() []byte {
	key := make([]byte, len(BrandClassKey))
	copy(key, BrandClassKey)
	return key
}
