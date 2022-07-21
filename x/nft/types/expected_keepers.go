package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/galaxies-labs/galaxy/x/brand/types"
)

type BrandKeeper interface {
	GetBrand(ctx sdk.Context, brandID string) (types.Brand, bool)
}
