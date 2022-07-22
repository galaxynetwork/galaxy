package types

import "strings"

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
	KeyPrefixClass = []byte{0x01}

	Delimiter = []byte{0x00}
)

func GetClassUniqueID(brandID, id string) string {
	return strings.Join([]string{brandID, id}, "/")
}

func GetClassStoreKey(brandID, id string) []byte {
	key := make([]byte, len(KeyPrefixClass)+len(brandID)+len(Delimiter)+len(id))
	copy(key, KeyPrefixClass)
	copy(key[len(KeyPrefixClass):], brandID)
	copy(key[len(KeyPrefixClass)+len(brandID):], Delimiter)
	copy(key[len(KeyPrefixClass)+len(brandID)+len(Delimiter):], id)
	return key
}

func GetClassOfBrandPrefix(brandID string) []byte {
	key := make([]byte, len(KeyPrefixClass)+len(brandID)+len(Delimiter))
	copy(key, KeyPrefixClass)
	copy(key[len(KeyPrefixClass):], brandID)
	copy(key[len(KeyPrefixClass)+len(brandID):], Delimiter)
	return key
}
