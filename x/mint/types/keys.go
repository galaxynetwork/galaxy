package types

var MinterKey = []byte{0x00}

const (
	// ModuleName defines the module name
	ModuleName = "mint"

	// StoreKey defines the primary module store key
	StoreKey    = ModuleName
	MemStoreKey = "mem" + ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	DefaultMintDenom = "uglx"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
