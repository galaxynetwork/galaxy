package types

var MinterKey = []byte{0x00}

const (
	// ModuleName defines the module name
	ModuleName = "clairdrop"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	DefaultClaimDenom = "uglx"

	ClaimRecordStorePrefix = "claim_recrod_store"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
