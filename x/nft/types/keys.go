package types

const (
	ModuleName = "nft"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

//(0x01|[]byte(brandID)|0x00|[]byte(classID) -> Class{})
// Query
//	- (0x01) prefix -> All classes
//	- (0x01|[]byte(brandID)) prefix -> All classes given brandID
//	- (0x01|[]byte(brandID)) prefix -> All classes given brandID

var (
	KeyPrefixClass = []byte{0x01}
)

func GetClassStoreKey(brandID, id string) []byte {
	key := make([]byte, len(KeyPrefixClass)+len(brandID)+len(id))
	copy(key, KeyPrefixClass)
	copy(key[len(KeyPrefixClass):], brandID)
	copy(key[len(KeyPrefixClass)+len(brandID):], id)
	return key
}

func GetClassOfBrandPrefix(brandID string) []byte {
	key := make([]byte, len(KeyPrefixClass)+len(brandID))
	copy(key, KeyPrefixClass)
	copy(key[len(KeyPrefixClass):], brandID)
	return key
}
