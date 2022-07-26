package types

import (
	"strconv"
	"strings"

	"github.com/galaxies-labs/galaxy/internal/conv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "nft"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

//0x01|[]byte(brandID)|0x00|[]byte(classID) -> Class{}
// Query
//	- (0x01) prefix -> All classes
//	- (0x01|[]byte(brandID)|0x00) prefix -> All classes given brandID
//(0x02|[]byte(brandID)|0x00|[]byte(classID) -> Supply{}
//
//0x03|[]byte(brandID)|0x00|[]byte(classID)|0x00|[]byte(nftID) -> NFT{}
//0x04|[]byte(brandID)|0x00|[]byte(classID)|0x00|[]byte(nftID) -> Owner{}

//authz
// brand
// nft
// marketplace

var (
	BrandClassKey       = []byte{0x01}
	BrandClassSupplyKey = []byte{0x02}
	NFTKey              = []byte{0x03}
	OwnerKey            = []byte{0x04}

	Delimiter = []byte{0x00}
)

func GetNFTUniqueID(brandID, classID string, id uint64) string {
	return strings.Join([]string{brandID, classID, strconv.FormatUint(id, 10)}, "/")
}

func ParseNFTUniqueID(uniqueID string) (string, string, uint64, error) {
	ids := strings.Split(uniqueID, "/")

	nftID, err := strconv.ParseUint(ids[2], 10, 64)
	if err != nil {
		return "", "", 0, err
	}

	return ids[0], ids[1], nftID, nil
}

func GetClassUniqueID(brandID, id string) string {
	return strings.Join([]string{brandID, id}, "/")
}

func ParseClassUniqueID(uniqueID string) (string, string) {
	ids := strings.Split(uniqueID, "/")
	return ids[0], ids[1]
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

func GetClassSupplyStoreKey(brandID string) []byte {
	brandIDBz := conv.UnsafeStrToBytes(brandID)

	key := make([]byte, len(BrandClassSupplyKey)+len(brandIDBz)+len(Delimiter))

	copy(key, BrandClassSupplyKey)
	copy(key[len(BrandClassSupplyKey):], brandIDBz)
	copy(key[len(BrandClassSupplyKey)+len(brandIDBz):], Delimiter)
	return key
}

func GetNFTStoreKey(brandID, classID string) []byte {
	brandIDBz := conv.UnsafeStrToBytes(brandID)
	classIDBz := conv.UnsafeStrToBytes(classID)

	key := make([]byte, len(NFTKey)+len(brandIDBz)+len(Delimiter)+len(classIDBz)+len(Delimiter))

	copy(key, NFTKey)
	copy(key[len(NFTKey):], brandIDBz)
	copy(key[len(NFTKey)+len(brandIDBz):], Delimiter)
	copy(key[len(NFTKey)+len(brandIDBz)+len(Delimiter):], classIDBz)
	copy(key[len(NFTKey)+len(brandIDBz)+len(Delimiter)+len(classIDBz):], Delimiter)
	return key
}

func GetOwnerStoreKey(brandID, classID string, id uint64) []byte {
	brandIDBz := conv.UnsafeStrToBytes(brandID)
	classIDBz := conv.UnsafeStrToBytes(classID)
	idBz := sdk.Uint64ToBigEndian(id)

	key := make([]byte, len(OwnerKey)+len(brandIDBz)+len(Delimiter)+len(classIDBz)+len(Delimiter)+len(idBz))

	copy(key, OwnerKey)
	copy(key[len(OwnerKey):], brandIDBz)
	copy(key[len(OwnerKey)+len(brandIDBz):], Delimiter)
	copy(key[len(OwnerKey)+len(brandIDBz)+len(Delimiter):], classIDBz)
	copy(key[len(OwnerKey)+len(brandIDBz)+len(Delimiter)+len(classIDBz):], Delimiter)
	copy(key[len(OwnerKey)+len(brandIDBz)+len(Delimiter)+len(classIDBz)+len(Delimiter):], idBz)
	return key
}
