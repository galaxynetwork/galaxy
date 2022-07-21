package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/nft/types"
)

type Keeper struct {
	storeKey sdk.StoreKey

	cdc         codec.BinaryCodec
	brandKeeper types.BrandKeeper
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, brandKeeper types.BrandKeeper) *Keeper {
	return &Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		brandKeeper: brandKeeper,
	}
}
